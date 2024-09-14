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
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		setErrorResponse(w, "invalid request", err)
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if signData.Password == storedPassword {
		jwtInstance := jwt.New(jwt.SigningMethodHS256)
		token, err := jwtInstance.SignedString([]byte(storedPassword))
		log.Println(token + " token из sign")
		taskData, err := json.Marshal(constans.SignInResponse{Token: token})
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(taskData)

		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}
	} else {
		errorResponse := constans.ErrorResponse{Error: "Неверный пароль"}
		errorData, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(errorData)

		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}
	}
}
