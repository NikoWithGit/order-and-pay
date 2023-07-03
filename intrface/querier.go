package intrface

import "database/sql"

type Querier interface {
	Query(s string, arg ...any) (*sql.Rows, error)
}
