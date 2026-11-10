package id

import (
	"strings"

	"github.com/api-go/jsmod"
)

func init() {
	jsmod.Register("id", map[string]any{
		"Make": jsMake,
	})
}

func jsMake(size int, forDB *string) string {
	dbType := ""
	if forDB != nil {
		dbType = strings.ToLower(*forDB)
	}

	switch dbType {
	case "mysql":
		return DefaultIDMaker.GetForMysql(size)
	case "postgres", "pg", "pgsql":
		return DefaultIDMaker.GetForPostgreSQL(size)
	default:
		return DefaultIDMaker.Get(size)
	}
}
