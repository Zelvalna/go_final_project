package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Zelvalna/go_final_project/config"

	"github.com/Zelvalna/go_final_project/model"

	"github.com/golang-jwt/jwt"
)

// SingInHandler проверка пароля и возврат JWT токена
func SingInHandler(w http.ResponseWriter, r *http.Request, cfg config.Config) {
	var signData model.SignInRequest
	// Декодируем запрос с паролем
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		setErrorResponse(w, "invalid request", err)
	}
	// Раз в auth изменили и тут изменим

	if cfg.TodoPassword != signData.Password || cfg.TodoPassword == "" {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(cfg.TodoPassword))
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
