package mysql

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

func Connection(dbConfig *config.Config) (*sql.DB, error) {
	if dbConfig == nil {
		dbConfig = config.Init()
	}

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		dbConfig.GetDBUser(),
		dbConfig.GetDBPassword(),
		dbConfig.GetDBHost(),
		dbConfig.GetDBPort(),
		dbConfig.GetDBName(),
	)
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbConn.SQL = db

	return dbConn.SQL, err
}
