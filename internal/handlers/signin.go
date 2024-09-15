package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zelvalna/go_final_project/constans"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
)

// SingInHandler проверка пароля и возврат JWT токена
func SingInHandler(w http.ResponseWriter, r *http.Request) {
	var signData constans.SignInRequest
	// Декодируем запрос с паролем
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		setErrorResponse(w, "invalid request", err)
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword != signData.Password || storedPassword == "" {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(storedPassword))
	log.Println(tokenString, err)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenString}); err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
	}
}
