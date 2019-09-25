package push

import (
	"time"

	"bitbucket.org/everymind/evmd-golib/logger"

	faktory "github.com/contribsys/faktory/client"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-gronos/model"
)

// Send envia um push de job para o Faktory
func Send(s model.JobScheduler) error {
	cl, err := faktory.Open()
	if err != nil {
		time.Sleep(5 * time.Second)
		cl, err = faktory.Open()
		if err != nil {
			return errors.Wrap(err, "faktory.Open()")
		}
	}

	params := []interface{}{
		s.ID,
		s.TenantID,
		s.TenantName,
		s.MiddlewareName,
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
		"dsn":   s.DSN,
		"stack": s.MiddlewareName,
	}

	if err = cl.Push(job); err != nil {
		return errors.Wrap(err, "cl.Push(job)")
	}

	logger.Tracef("Job '%s', Function: '%s', Stack: '%s', Params: '%v' pushed to Faktory on queue '%s'", s.JobName, s.FuncName, s.MiddlewareName, params, s.Queue)

	return nil
}
