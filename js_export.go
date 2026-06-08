package id

import (
	"strings"

	"apigo.cc/go/jsmod"
)

func init() {
	jsmod.Register("id", map[string]any{
		"Make": func(size int, forDB ...string) string {
			dbType := ""
			if len(forDB) > 0 {
				dbType = strings.ToLower(forDB[0])
			}

			switch dbType {
			case "mysql":
				return DefaultIDMaker.GetForMysql(size)
			case "postgres", "pg", "pgsql":
				return DefaultIDMaker.GetForPostgreSQL(size)
			default:
				return DefaultIDMaker.Get(size)
			}
		},
	})
}
