package db

import (
	"database/sql"
	"order-and-pay/iface"
)

type SqlTx struct {
	*sql.Tx
}

func NewSqlTx(tx *sql.Tx) iface.Itx {
	return &SqlTx{tx}
}

func (tx *SqlTx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *SqlTx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *SqlTx) Begin() (iface.Itx, error) {
	return nil, nil
}
