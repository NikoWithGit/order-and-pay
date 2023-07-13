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

func (db *SqlDb) Begin() (intrface.Itx, error) {
	return db.DB.Begin()
}
