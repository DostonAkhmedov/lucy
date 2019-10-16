package mysql

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/config"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func Connection(conf *config.Config) (*sql.DB, error) {
	if conf == nil {
		conf = config.Init()
	}

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.GetDBUser(),
		conf.GetDBPassword(),
		conf.GetDBHost(),
		conf.GetDBPort(),
		conf.GetDBName(),
	)
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dbConn.SQL = db

	return dbConn.SQL, err
}
