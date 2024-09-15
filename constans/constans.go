package constans

import "github.com/golang-jwt/jwt"

const (
	DefPort = "7540"
	WebDir  = "./web"
	DatePat = "20060102"
)

type Task struct {
	ID      string `json:"id,omitempty" db:"id"`
	Date    string `json:"date,omitempty" db:"date"`
	Title   string `json:"title,omitempty" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat,omitempty" db:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskIdResponse struct {
	Id int `json:"id"`
}
type Tasks struct {
	Tasks []Task `json:"tasks"`
}
type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
}
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.StandardClaims
}
