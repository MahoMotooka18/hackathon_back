package controller

import (
	"encoding/json"
	"github.com/oklog/ulid"
	"golang.org/x/crypto/bcrypt"
	"hackathon/dao"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

func generateULID() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	return ulid.MustNew(ms, entropy)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "エラー", err
	}
	return string(hashedPassword), nil
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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

	if requestBody.Name == "" {
		log.Printf("名前が入力されていません , \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(requestBody.Name) > 50 {
		log.Printf("名前の文字数が50を超えています , \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Email == "" {
		log.Printf("Eメールアドレスが入力されていません, \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Password == "" {
		log.Printf("パスワードが入力されていません, \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pattern := "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[^A-Za-z\\d\\s]).+$"
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(requestBody.Password) {
		log.Printf("パスワードは英数字と特殊記号を含む必要があります\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Hashed, err := hashPassword(requestBody.Password)
	if err != nil {
		log.Printf("パスワードのハッシュ化に失敗: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ulidId := generateULID()
	tx, err := dao.DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Printf("transactionの開始に失敗しました , %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO user(id, name, email, password) VALUES (?,?,?,?)", ulidId, requestBody.Name, requestBody.Email, Hashed)
	if err != nil {
		tx.Rollback()
		log.Printf("データの登録に失敗しました, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("コミットに失敗しました , %v\n", err)
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
