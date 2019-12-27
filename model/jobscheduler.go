package model

import (
	"database/sql"
)

type JobScheduler struct {
	ID                int64          `db:"id"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	StackID           int16          `db:"stack_id"`
	StackName         string         `db:"stack_name"`
	DSN               string         `db:"dsn"`
	JobName           string         `db:"job_name"`
	FuncName          string         `db:"function_name"`
	Queue             string         `db:"queue"`
	Cron              string         `db:"cron"`
	Parameters        sql.NullString `db:"parameters"` // In DB is JSONB
	Retry             int16          `db:"retry"`
	AllowsConcurrency bool           `db:"allows_concurrency"`
	AllowsSchedule    bool           `db:"allows_schedule"`
	ScheduleTime      int16          `db:"schedule_time"`
	IsActive          bool           `db:"is_active"`
	IsDeleted         bool           `db:"is_deleted"`
	OrgID             string         `db:"org_id"`
	AppEngineName     string         `db:"appengine_name"`
}

func (s JobScheduler) Copy() JobScheduler {
	return JobScheduler{
		ID:                s.ID,
		TenantID:          s.TenantID,
		TenantName:        s.TenantName,
		StackID:           s.StackID,
		StackName:         s.StackName,
		DSN:               s.DSN,
		JobName:           s.JobName,
		FuncName:          s.FuncName,
		Queue:             s.Queue,
		Cron:              s.Cron,
		Parameters:        s.Parameters,
		Retry:             s.Retry,
		AllowsConcurrency: s.AllowsConcurrency,
		AllowsSchedule:    s.AllowsSchedule,
		ScheduleTime:      s.ScheduleTime,
		IsActive:          s.IsActive,
		IsDeleted:         s.IsDeleted,
		OrgID:             s.OrgID,
		AppEngineName:     s.AppEngineName,
	}
}
