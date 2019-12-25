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
