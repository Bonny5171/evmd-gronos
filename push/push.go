package push

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/everymind/evmd-golib/db"

	"bitbucket.org/everymind/evmd-golib/db/dao"
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

	if err = cleanGhostJobs(s.DSN); err != nil {
		return fmt.Errorf("cleanGhostJobs(db): %v", err)
	}

	if err = cl.Push(job); err != nil {
		return fmt.Errorf("cl.Push(job): %w", err)
	}

	logger.Tracef("Job '%s', Function: '%s', Stack: '%s', Params: '%v' pushed to Faktory on queue '%s'", s.JobName, s.FuncName, s.StackName, params, s.Queue)

	//pingJob(s.AppEngineName.String)

	return nil
}

func cleanGhostJobs(dsn string) (err error) {
	if err := db.Create(&db.PostgresDB{
		ConnectionStr: dsn,
		MaxLifetime:   0,
		MaxIdleConns:  1,
		MaxOpenConns:  1,
	}); err != nil {
		return fmt.Errorf("db.Create(): %v", err)
	}

	dbConn := db.Conn
	defer dbConn.Close()

	logger.Tracef("Calling CleanGhostJob on: %v", dbConn)
	if err = dao.CleanGhostJobs(dbConn); err != nil {
		return err
	}
	logger.Tracef("CleanGhostJob complete.")
	return nil
}
