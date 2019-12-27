package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/everymind/evmd-golib/logger"
	"bitbucket.org/everymind/evmd-golib/utils"
	"github.com/besser/cron"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"

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

			entryName := fmt.Sprintf("%03d_%s_%s_%s", j.TenantID, utils.RemoveSpacesAndLower(j.StackName), utils.RemoveSpacesAndLower(j.JobName), utils.RemoveSpacesAndLower(j.FuncName))
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
				appName := j.AppEngineName

				// Anonymous function
				fn := func() {
					pingJob(appName)

					if err := push.Send(s); err != nil {
						logger.Errorln(errors.Wrap(err, "push.Send()"))
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

func pingJob(appEngineName string) {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		logger.Errorln(errors.Wrap(err, "google.FindDefaultCredentials()"))
	}

	var sb strings.Builder
	sb.WriteString("https://")
	sb.WriteString(appEngineName)
	sb.WriteString("-dot-")
	sb.WriteString(credentials.ProjectID)
	sb.WriteString(".appspot.com/_ah/start")

	response, err := http.Get(sb.String())
	if err != nil {
		logger.Errorln(errors.Wrap(err, "http.Get()"))
	}

	data, _ := ioutil.ReadAll(response.Body)
	logger.Infoln("ping job: " + string(data))
}
