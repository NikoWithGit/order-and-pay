package iface

import "database/sql"

type Idb interface {
	Begin() (Itx, error)
	Query(string, ...any) (*sql.Rows, error)
}

type Itx interface {
	Idb
	Rollback() error
	Commit() error
}
