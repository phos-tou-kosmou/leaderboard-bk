package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

//var jwtKey = []byte("my_secret_key")
//
//var users = map[string]string{
//	"user1": "password1",
//	"user2": "password2",
//}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct to read the username and password from the request body
type User struct {
	*Credentials
	CreatedAt time.Time
}
// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Student_test struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	GPA int `json:"gpa"`
	Sport string `json:"sport"`
}

type Student struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	GPA float32 `json:"gpa"`
	Sport string `json:"sport"`
	CreatedAt time.Time `json:"t_stamp"`
}