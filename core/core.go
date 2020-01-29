package core

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3" // "github.com/besser/cron"
	"github.com/spf13/cast"

	"bitbucket.org/everymind/evmd-golib/logger"
	"bitbucket.org/everymind/evmd-golib/utils"
	"bitbucket.org/everymind/evmd-gronos/dao"
	"bitbucket.org/everymind/evmd-gronos/push"
)

type ScheduledJob struct {
	ID   cron.EntryID
	Spec string
}

// Run é onde se inicia o processo
func Run(c *cron.Cron, sJobs map[string]ScheduledJob) error {
	tenantID := cast.ToInt(os.Getenv("TENANT_ID"))

	// Recupera todas os 'jobs' que deverão ser executados
	jobs, err := dao.GetSchedules(tenantID)
	if err != nil {
		return fmt.Errorf("dao.GetSchedules(): %w", err)
	}

	if len(jobs) > 0 {
		for _, j := range jobs {
			var insert bool

			entryName := fmt.Sprintf("%03d_%s_%s_%s", j.TenantID, utils.RemoveSpacesAndLower(j.StackName), utils.RemoveSpacesAndLower(j.JobName), utils.RemoveSpacesAndLower(j.FuncName))

			sJob, ok := sJobs[entryName]
			if ok && sJob.ID > 0 {
				if entry := c.Entry(sJob.ID); entry.ID > 0 {
					if !j.Cron.Valid || j.Cron.String != sJob.Spec || !j.IsActive || j.IsDeleted {
						logger.Infof("Removing scheduled job '%s(id:%d)'", entryName, entry.ID)
						c.Remove(entry.ID)
						if !j.IsDeleted {
							insert = false
						} else {
							insert = j.IsActive
						}
					}
				} else if j.Cron.Valid && j.IsActive && !j.IsDeleted {
					insert = true
				}
			} else if j.Cron.Valid && j.IsActive && !j.IsDeleted {
				insert = true
			}

			if insert {
				s := j.Copy()

				if j.AppEngineName.Valid == false {
					err = errors.New("AppEngineName not found")
					return err
				}

				// Anonymous function
				fn := func() {
					if err := push.Send(s); err != nil {
						logger.Errorln(fmt.Errorf("push.Send(): %w", err))
					}

					pingJob(j.AppEngineName.String)
				}

				location, _ := dao.GetParamByOrgID(j.OrgID, "ORG_TZ_LOCATION")
				cronSpec := j.Cron.String

				if location != "UTC" {
					_, err := time.LoadLocation(location)
					if err != nil {
						logger.Errorln(fmt.Errorf("time.LoadLocation(): %w", err))
						logger.Warningf("The location '%s' is invalid, setting to UTC", location)
					} else {
						cronSpec = fmt.Sprintf("CRON_TZ=%s %s", location, j.Cron.String)
					}
				}

				id, err := c.AddFunc(cronSpec, fn)
				if err != nil {
					return err
				}

				sJobs[entryName] = ScheduledJob{
					ID:   id,
					Spec: j.Cron.String,
				}

				entry := c.Entry(id)
				logger.Infof("Next scheduling job '%s' (id:%d, cron:'%s') to: %v", entryName, entry.ID, j.Cron.String, entry.Next)
			}
		}
	} else {
		// Removendo todos os jobs agendados
		for _, e := range c.Entries() {
			c.Remove(e.ID)
		}

		logger.Infoln("No scheduled jobs!")
	}

	return nil
}

func pingJob(appEngineName string) {
	cloudProject := os.Getenv("GOOGLE_CLOUD_PROJECT")

	var sb strings.Builder
	sb.WriteString("https://")
	sb.WriteString(appEngineName)
	sb.WriteString("-dot-")
	sb.WriteString(cloudProject)
	sb.WriteString(".appspot.com/_ah/start")

	response, err := http.Get(sb.String())
	if err != nil {
		logger.Errorln(fmt.Errorf("http.Get(): %w", err))
	}

	logger.Infof("ping to job '%s' at %s: %s", appEngineName, sb.String(), response.Status)

	if response.StatusCode/100 != 2 {
		err := fmt.Errorf("job %s unavaliable", appEngineName)
		logger.Errorln(err)
	}
}
