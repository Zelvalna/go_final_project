package middleware

import (
	"log"
	"net/http"

	"github.com/Zelvalna/go_final_project/config"

	"github.com/golang-jwt/jwt"
)

// Auth проверяет JWT токен
func Auth(nextHandler http.HandlerFunc, cfg config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//получаем переменную окружения
		if len(cfg.TodoPassword) > 0 {

			// получаем куку с токеном
			var tokenString string
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString = cookie.Value
			jwtInstance := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtInstance.SignedString([]byte(cfg.TodoPassword))
			log.Println(token + " token из AUTH")
			if tokenString != token {
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}
		}
		nextHandler(w, r)
	})
}
