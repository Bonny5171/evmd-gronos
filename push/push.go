package push

import (
	"fmt"
	"time"

	faktory "github.com/contribsys/faktory/client"

	"bitbucket.org/everymind/evmd-golib/v2/logger"
	"bitbucket.org/everymind/evmd-gronos/v3/model"
)

// Send envia um push de job para o Faktory
func Send(s model.JobScheduler) error {
	cl, err := faktory.Open()
	if err != nil {
		time.Sleep(5 * time.Second)
		cl, err = faktory.Open()
		if err != nil {
			return fmt.Errorf("faktory.Open(): %w", err)
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

	return nil
}
