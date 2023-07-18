package db

import (
	"database/sql"
	"order-and-pay/intrface"
)

type SqlTx struct {
	*sql.Tx
}

func NewSqlTx(tx *sql.Tx) intrface.Itx {
	return &SqlTx{tx}
}

func (tx *SqlTx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *SqlTx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *SqlTx) Begin() (intrface.Itx, error) {
	return nil, nil
}
