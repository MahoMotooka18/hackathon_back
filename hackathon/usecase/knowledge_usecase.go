package usecase

import (
	"hackathon/knowledge_controller"
	"log"
	"net/http"
)

func KnowledgeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8050")
	switch r.Method {
	case http.MethodGet:
		knowledge_controller.KnowlegdeGetHandler(w, r)

	case http.MethodPost:
		knowledge_controller.KnowlegdePostHandler(w, r)

	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
