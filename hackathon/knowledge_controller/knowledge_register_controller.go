package knowledge_controller

import (
	"encoding/json"
	"github.com/oklog/ulid"
	"hackathon/dao"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func generateULID() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	return ulid.MustNew(ms, entropy)
}

func KnowlegdePostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name       string `json:"name"`
		Url        string `json:"url"`
		Category   string `json:"category"`
		Curriculum string `json:"curriculum"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("fail:request body is empty , %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" {
		log.Printf("エラー: 名前が入力されていません , \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(requestBody.Name) > 50 {
		log.Printf("名前の文字数が50を超えています , \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Url == "" {
		log.Printf("urlが入力されていません, \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Category == "" {
		log.Printf("カテゴリーが選択されていません, \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Curriculum == "" {
		log.Printf("カリキュラムが選択されていません, \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ulidId := generateULID()
	tx, err := dao.DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Printf("情報の取得に失敗しました, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO knowledge(id, name, url, category, curriculum) VALUES (?,?,?,?,?)", ulidId, requestBody.Name, requestBody.Url, requestBody.Category, requestBody.Curriculum)
	if err != nil {
		tx.Rollback()
		log.Printf("情報の登録に失敗しました , %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("fail:コミットに失敗しました , %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]string{"id": ulidId.String()}
	jsonResponse, err := json.Marshal(responseJSON)
	if err != nil {
		log.Printf("Marshalingに失敗しました , %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
