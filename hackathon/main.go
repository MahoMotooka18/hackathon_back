package main

import (
	"database/sql"
	"fmt"
	"hackathon/dao"
	"hackathon/knowledge_controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var db *sql.DB

func Init() {
	//mysqlUser := os.Getenv("MYSQL_USER")
	//mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	//mysqlHost := os.Getenv("MYSQL_HOST")
	//mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	_db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	db = _db
}

func main() {
	//  /knowledgeでリクエストされたら
	http.HandleFunc("/knowledge", knowledge_controller.KnowlegdeGetHandler)
	http.HandleFunc("/knowledgepost", knowledge_controller.KnowlegdePostHandler)
	// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
	closeDBWithSysCall()
	// 8050番ポートでリクエストを待ち受ける
	log.Println("Listening...")
	if err := http.ListenAndServe(":8050", nil); err != nil {
		log.Fatal(err)
	}
}

// ③ Ctrl+CでHTTPサーバー停止時にDBをクローズする
func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)
		if err := dao.DB.Close(); err != nil {
			log.Fatal(err)
			//fmt.Println("cannot close sql")
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
