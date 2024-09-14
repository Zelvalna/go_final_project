package middleware

import (
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
)

// Auth проверяет JWT токен
func Auth(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//получаем переменную окружения
		storedPassword := os.Getenv("TODO_PASSWORD")
		if len(storedPassword) > 0 {

			// получаем куку с токеном
			var tokenString string
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			jwtInstance := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtInstance.SignedString([]byte(storedPassword))
			log.Println(token + " token из AUTH")
			if tokenString != token {
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
			}
		}
		nextHandler(w, r)
	}
}
