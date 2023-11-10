package usecase

import (
	"hackathon/user_controller"
	"log"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8050")
	switch r.Method {
	case http.MethodGet:
		user_controller.UserGetHandler(w, r)

	case http.MethodPost:
		user_controller.UserPostHandler(w, r)

	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
