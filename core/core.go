package core

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/everymind/evmd-golib/logger"
	"bitbucket.org/everymind/evmd-golib/utils"
	"github.com/besser/cron"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"bitbucket.org/everymind/evmd-gronos/dao"
	"bitbucket.org/everymind/evmd-gronos/push"
)

// Run é onde se inicia o processo
func Run(c *cron.Cron) error {
	tenantID := cast.ToInt(os.Getenv("TENANT_ID"))

	// Recupera todas os 'jobs' que deverão ser executados
	jobs, err := dao.GetSchedules(tenantID)
	if err != nil {
		return errors.Wrap(err, "dao.GetSchedules()")
	}

	if len(jobs) > 0 {
		for _, j := range jobs {
			var insert bool

			entryName := fmt.Sprintf("%s_%s_%s_%s", utils.RemoveSpacesAndLower(j.TenantName), utils.RemoveSpacesAndLower(j.StackName), utils.RemoveSpacesAndLower(j.JobName), utils.RemoveSpacesAndLower(j.FuncName))
			entry := c.EntryName(entryName)

			if entry.ID > 0 {
				if j.Cron != entry.Spec || !j.IsActive || j.IsDeleted {
					logger.Infof("Removing scheduled job '%s'", entry.Name)
					c.Remove(entry.ID)
					if !j.IsDeleted {
						insert = false
					} else {
						insert = j.IsActive
					}
				}
			} else if j.IsActive && !j.IsDeleted {
				insert = true
			}

			if insert {
				s := j.Copy()

				// Anonymous function
				fn := func() {
					if err := push.Send(s); err != nil {
						logger.Errorln(errors.Wrap(err, "push.LoadAndPush()"))
					}
				}

				location, _ := dao.GetParamByOrgID(j.OrgID, "ORG_TZ_LOCATION")

				var id cron.EntryID
				if location == "UTC" {
					id, err = c.AddFuncN(entryName, j.Cron, fn)
					if err != nil {
						return err
					}
				} else {
					loc, err := time.LoadLocation(location)
					if err != nil {
						logger.Errorln(errors.Wrap(err, "time.LoadLocation()"))
						logger.Warningf("The location '%s' is invalid, setting to UTC", location)
						loc = time.UTC
					}
					id, err = c.AddFuncNLocation(entryName, j.Cron, loc, fn)
					if err != nil {
						return err
					}
				}

				entry = c.Entry(id)
				logger.Infof("Next scheduling job '%s' ('%s') to: %v", entry.Name, entry.Spec, entry.Next)
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
