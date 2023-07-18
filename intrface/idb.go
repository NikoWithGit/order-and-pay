package intrface

import "database/sql"

type Idb interface {
	Itx
	Begin() (Idb, error)
	Query(string, ...any) (*sql.Rows, error)
}

type Itx interface {
	Rollback() error
	Commit() error
}
