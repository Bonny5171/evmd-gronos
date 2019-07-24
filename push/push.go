package push

import (
	"os"
	"time"

	"bitbucket.org/everymind/gopkgs/logger"

	faktory "github.com/contribsys/faktory/client"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gronos/model"
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
		os.Getenv("GRONOS_STACK_NAME"),
		s.AllowsConcurrency,
		s.AllowsSchedule,
		s.ScheduleTime,
	}

	if s.Parameters.Valid {
		params = append(params, s.Parameters.String)
	}

	job := faktory.NewJob(s.JobName, params...)
	job.Queue = s.Queue
	job.Retry = int(s.Retry)

	job.Custom = map[string]interface{}{
		"dsn":   os.Getenv("GRONOS_DATABASE_DSN"),
		"stack": os.Getenv("GRONOS_STACK_NAME"),
	}

	if err = cl.Push(job); err != nil {
		return errors.Wrap(err, "cl.Push(job)")
	}

	logger.Tracef("Job '%s(%v)' pushed to Faktory on queue '%s'", s.JobName, params, s.Queue)

	return nil
}
