package dao

import (
	"fmt"

	"bitbucket.org/everymind/evmd-golib/v2/db"
	dd "bitbucket.org/everymind/evmd-golib/v2/db/dao"
)

func GetParamByOrgID(orgID, paramName string) (string, error) {
	conn, err := db.GetConnection("CONFIG")
	if err != nil {
		return "", fmt.Errorf("db.GetConnection('CONFIG'): %w", err)
	}

	p, err := dd.GetParameterByOrgID(conn, orgID, paramName)
	if err != nil {
		return "", err
	}

	return p.Value, nil
}
