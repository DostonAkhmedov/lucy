package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBLogin, DBPassword, DBHost, DBPort, DBName);
	DB, err = sql.Open(DBType, connection)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
