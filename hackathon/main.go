package main

import (
	"hackathon/dao"
	"hackathon/knowledge_controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
