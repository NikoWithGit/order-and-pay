package db

import (
	"database/sql"
	"order-and-pay/env"
	"order-and-pay/iface"
)

type SqlDb struct {
	*sql.DB
}

func NewSqlDb() (*SqlDb, error) {
	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	return &SqlDb{db}, err
}

func (db *SqlDb) Begin() (iface.Itx, error) {
	tx, err := db.DB.Begin()
	return NewSqlTx(tx), err
}
