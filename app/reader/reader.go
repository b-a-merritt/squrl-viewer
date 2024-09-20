package reader

import (
	"github.com/b-a-merritt/squrlviewer/app/config"

	"fmt"
)

func Init(cfg *config.Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("")
	}

	user := ""
	password := "password="
	dbname := ""
	schema := "schema=public"
	sslMode := "sslmode=verify-full"

	cfg.ConnStr = fmt.Sprintf(`%v %v %v %v %v`, user, password, dbname, schema, sslMode)

	cfg.ConnStr = "postgres://platform:laksjdflcnwkl2342@localhost:5432/soulrefiner?sslmode=disable"

	return nil
}
