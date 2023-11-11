package main

import (
	"hackathon/dao"
	"hackathon/usecase"
	"hackathon/user_controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// データベース接続の初期化
	dao.Init()

	//  /Signupでリクエストされたら
	http.HandleFunc("/signup", user_controller.UserPostHandler)

	http.HandleFunc("/login", user_controller.UserGetHandler)
	//  /knowledgeでリクエストされたら
	http.HandleFunc("/knowledge", usecase.KnowledgeHandler)
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
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
