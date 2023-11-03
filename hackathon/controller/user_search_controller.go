package controller

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"hackathon/dao"
	"hackathon/model"
	"log"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// ユーザーの認証情報をリクエストから受け取る
	var requestBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("リクエストボディが空です , %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// データベースからユーザーのハッシュ化されたパスワードを取得
	var hashedPassword string
	if err := dao.DB.QueryRow("SELECT password FROM user WHERE name=?", requestBody.Name).Scan(&hashedPassword); err != nil {
		log.Printf("ユーザーが見つかりません: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// パスワードの一致を確認
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestBody.Password)); err != nil {
		log.Printf("パスワードが一致しません: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// ユーザー情報を取得してクライアントに返す
	user := model.UserResForHTTPGet{}
	if err := dao.DB.QueryRow("SELECT id, name, email FROM user WHERE name=?", requestBody.Name).Scan(&user.Id, &user.Name, &user.Email); err != nil {
		log.Printf("ユーザー情報の取得に失敗しました: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	bytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("JSONのエンコードに失敗しました: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
