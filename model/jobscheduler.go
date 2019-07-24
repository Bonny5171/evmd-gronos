package model

import (
	"database/sql"
)

type JobScheduler struct {
	ID                int64          `db:"id"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	Name              string         `db:"name"`
	Queue             string         `db:"queue"`
	JobName           string         `db:"job_name"`
	Parameters        sql.NullString `db:"parameters"` // In DB is JSONB
	Retry             int16          `db:"retry"`
	Description       sql.NullString `db:"description"`
	Cron              string         `db:"cron"`
	OrgID             string         `db:"org_id"`
	AllowsConcurrency bool           `db:"allows_concurrency"`
	AllowsSchedule    bool           `db:"allows_schedule"`
	ScheduleTime      int16          `db:"schedule_time"`
	DocMetaData       sql.NullString `db:"doc_meta_data"` // In DB is JSONB
	IsActive          bool           `db:"is_active"`
	IsDeleted         bool           `db:"is_deleted"`
}

func (s JobScheduler) Copy() JobScheduler {
	return JobScheduler{
		ID:                s.ID,
		TenantID:          s.TenantID,
		TenantName:        s.TenantName,
		Name:              s.Name,
		Queue:             s.Queue,
		JobName:           s.JobName,
		Parameters:        s.Parameters,
		Retry:             s.Retry,
		Description:       s.Description,
		Cron:              s.Cron,
		OrgID:             s.OrgID,
		AllowsConcurrency: s.AllowsConcurrency,
		AllowsSchedule:    s.AllowsSchedule,
		ScheduleTime:      s.ScheduleTime,
		DocMetaData:       s.DocMetaData,
		IsActive:          s.IsActive,
		IsDeleted:         s.IsDeleted,
	}
}
