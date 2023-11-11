package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv"
	"log"
	"os"
)

var DB *sql.DB

func Init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	log.Printf("MYSQL_USER: %s", mysqlUser)
	log.Printf("MYSQL_PASSWORD: %s", mysqlPwd)
	log.Printf("MYSQL_HOST: %s", mysqlHost)
	log.Printf("MYSQL_DATABASE: %s", mysqlDatabase)

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	_db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	DB = _db
}
