package db

import (
	"database/sql"
	"order-and-pay/env"
	"order-and-pay/intrface"
)

type SqlDb struct {
	*sql.DB
}

func NewSqlDb() (*SqlDb, error) {
	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	return &SqlDb{db}, err
}

func (db *SqlDb) Begin() (intrface.Idb, error) {
	tx, err := db.DB.Begin()
	return NewSqlTx(tx), err
}

func (db *SqlDb) Rollback() error {
	return nil
}

func (db *SqlDb) Commit() error {
	return nil
}
