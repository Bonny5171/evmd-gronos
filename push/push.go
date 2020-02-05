package push

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/everymind/evmd-golib/logger"
	"bitbucket.org/everymind/evmd-gronos/model"
	faktory "github.com/contribsys/faktory/client"
	"github.com/spf13/cast"
)

// Send envia um push de job para o Faktory
func Send(s model.JobScheduler) error {
	cl, err := faktory.Open()
	if err != nil {
		logger.Errorln(fmt.Errorf("faktory.Open(): %w", err))

		count := cast.ToInt(os.Getenv("FAKTORY_CONNECTION_TRIES"))
		if count <= 0 {
			count = 10
		}

		for i := 0; i < count; i++ {
			time.Sleep(5 * time.Second)
			cl, err = faktory.Open()
			if err != nil {
				logger.Errorln(fmt.Errorf("faktory.Open(): %w", err))
				continue
			}
			break
		}
	}

	params := []interface{}{
		s.ID,
		s.JobName,
		s.TenantID,
		s.TenantName,
		s.StackName,
		s.AllowsConcurrency,
		s.AllowsSchedule,
		s.ScheduleTime,
	}

	if s.Parameters.Valid {
		params = append(params, s.Parameters.String)
	}

	job := faktory.NewJob(s.FuncName, params...)
	job.Queue = s.Queue
	job.Retry = int(s.Retry)

	job.Custom = map[string]interface{}{
		"dsn":   s.DSN + fmt.Sprintf(" application_name='%s'", s.JobName),
		"stack": s.StackName,
	}

	if err = cl.Push(job); err != nil {
		return fmt.Errorf("cl.Push(job): %w", err)
	}

	logger.Tracef("Job '%s', Function: '%s', Stack: '%s', Params: '%v' pushed to Faktory on queue '%s'", s.JobName, s.FuncName, s.StackName, params, s.Queue)

	pingJob(s.AppEngineName.String)

	return nil
}

func pingJob(appEngineName string) {
	cloudProject := os.Getenv("GOOGLE_CLOUD_PROJECT")

	var sb strings.Builder
	sb.WriteString("http://")
	sb.WriteString(appEngineName)
	sb.WriteString("-dot-")
	sb.WriteString(cloudProject)
	sb.WriteString(".appspot.com/")

	response, err := http.Get(sb.String())
	if err != nil {
		logger.Errorln(fmt.Errorf("http.Get(): %w", err))
	}

	if response != nil {
		logger.Infof("ping to job '%s' at %s: %s", appEngineName, sb.String(), response.Status)

		if response.StatusCode/100 != 2 {
			err := fmt.Errorf("job %s unavaliable", appEngineName)
			logger.Errorln(err)
		}
	}
}
