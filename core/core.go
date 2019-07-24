package core

import (
	"time"

	"bitbucket.org/everymind/gopkgs/logger"
	"bitbucket.org/everymind/gopkgs/utils"
	"github.com/besser/cron"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gronos/dao"
	"bitbucket.org/everymind/gronos/push"
)

// Run é onde se inicia o processo
func Run(c *cron.Cron) error {
	// Recupera todas os 'jobs' que deverão ser executados
	jobs, err := dao.GetSchedules()
	if err != nil {
		return errors.Wrap(err, "dao.GetSchedules()")
	}

	if len(jobs) > 0 {
		for _, j := range jobs {
			var insert bool

			entryName := utils.RemoveSpacesAndLower(j.TenantName) + "_" + j.Name
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
