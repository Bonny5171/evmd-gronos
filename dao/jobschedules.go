package dao

import (
	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-gronos/model"
)

// GetSchedules retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedules() (s []model.JobScheduler, err error) {
	conn, err := db.GetConnection("CONFIG")
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConnection('CONFIG')")
	}

	const query = `
	  SELECT j.id, 
	         j.tenant_id, 
	         t."name" AS tenant_name, 
	         j."name", 
	         j.queue, 
	         j.job_name, 
	         j.parameters, 
	         j.retry, 
	         j.description, 
	         j.cron, 
			 j.allows_concurrency, 
			 j.allows_schedule, 
			 j.schedule_time, 
	         j.is_active, 
	         j.is_deleted, 
	         t.org_id
	    FROM itgr.job_scheduler j
	   INNER JOIN public.tenant t ON j.tenant_id = t.id
	   ORDER BY j.id;`

	err = conn.Select(&s, query)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return s, nil
}
