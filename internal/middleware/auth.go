package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// Auth проверяет JWT токен
func Auth(nextHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//получаем переменную окружения
		storedPassword := os.Getenv("TODO_PASSWORD")
		if len(storedPassword) > 0 {

			// получаем куку с токеном
			var tokenString string
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString = cookie.Value
			jwtInstance := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtInstance.SignedString([]byte(storedPassword))
			log.Println(token + " token из AUTH")
			if tokenString != token {
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}
		}
		nextHandler(w, r)
	})
}
