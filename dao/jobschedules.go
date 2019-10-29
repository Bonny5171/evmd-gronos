package dao

import (
	"os"
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"bitbucket.org/everymind/evmd-gronos/model"
)

// GetSchedules retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedules(tenantID int) (s []model.JobScheduler, err error) {
	conn, err := db.GetConnection("CONFIG")
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConnection('CONFIG')")
	}

	tenantType := "JOB"
	if cast.ToBool(os.Getenv("DEBUG")) {
		tenantType = "DEBUG"
	}

	params := []interface{}{tenantType}
	query := strings.Builder{}

	query.WriteString(`
		SELECT j.id, 
			j.tenant_id, 
			t."name" AS tenant_name, 
			j.stack_id,
			m."name" AS stack_name,
			d.string AS dsn,
			j.job_name,
			j.function_name,
			j.queue,
			j.cron,   
			j.parameters, 
			j.retry, 
			j.allows_concurrency, 
			j.allows_schedule, 
			j.schedule_time, 
			j.is_active, 
			j.is_deleted, 
			t.org_id
		FROM public.job_scheduler j
		INNER JOIN public.tenant t ON j.tenant_id = t.id
		INNER JOIN public.stack  m ON j.stack_id = m.id AND m.is_active = TRUE AND m.is_deleted = FALSE
		INNER JOIN public.dsn    d ON m.id = d.stack_id AND upper(d."type") = $1`)

	if tenantID > 0 {
		query.WriteString(" WHERE j.tenant_id = $2")
		params = append(params, tenantID)
	}

	query.WriteString(" ORDER BY j.id;")

	err = conn.Select(&s, query.String(), params...)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return s, nil
}
