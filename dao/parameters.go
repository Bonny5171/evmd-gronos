package dao

import (
	"bitbucket.org/everymind/evmd-golib/db"
	dd "bitbucket.org/everymind/evmd-golib/db/dao"
	"github.com/pkg/errors"
)

func GetParamByOrgID(orgID, paramName string) (string, error) {
	conn, err := db.GetConnection("CONFIG")
	if err != nil {
		return "", errors.Wrap(err, "db.GetConnection('CONFIG')")
	}

	p, err := dd.GetParameterByOrgID(conn, orgID, paramName)
	if err != nil {
		return "", err
	}

	return p.Value, nil
}
