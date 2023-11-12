package knowledge_controller

import (
	"encoding/json"
	"hackathon/dao"
	"hackathon/model"
	"log"
	"net/http"
)

func KnowlegdeGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	rows, err := dao.DB.Query("SELECT id, name, url, category, curriculum FROM knowledge")
	if err != nil {
		log.Printf("fail: データを持ってくるのに失敗, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	knowledge := make([]model.KnowledgeResForHTTPGet, 0)
	for rows.Next() {
		var u model.KnowledgeResForHTTPGet
		if err := rows.Scan(&u.Id, &u.Name, &u.Url, &u.Category, &u.Curriculum); err != nil {
			log.Printf("fail: rows.Scan, %v\n", err)
			if err := rows.Close(); err != nil { // 500を返して終了するが、その前にrowsのClose処理が必要
				log.Printf("fail: rows.Close(), %v\n", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		knowledge = append(knowledge, u)
	}
	// ②-4
	bytes, err := json.Marshal(knowledge)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)

}
