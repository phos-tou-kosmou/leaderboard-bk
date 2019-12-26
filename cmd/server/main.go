package main

import (
	"encoding/json"
	"fmt"
	"leaderboard-bk/cmd/controllers"
	"leaderboard-bk/cmd/models"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

/*******************STUDENT API ROUTES****************************/
type test_struct struct {FirstName string `json:"first_name"`}
func StudentsIndex(w http.ResponseWriter, r *http.Request) {controllers.IndexStudents(w, r)}
func CreateStudent(w http.ResponseWriter, r *http.Request) {controllers.InsertStudent(w, r)}
func FetchStudent(w http.ResponseWriter, r *http.Request) {controllers.IndexStudents(w, r)}
//func UpdateStudent(w http.ResponseWriter, r *http.Request) {controllers.EditOfStudent(w, r)}
func DeleteStudent(w http.ResponseWriter, r *http.Request) {controllers.DeleteStudent(w, r)}
/*****************************************************************/


/*******************SIGN IN AND LANDING**************************/
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//authy := queryAuth(creds)
	//
	//if authy {
	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(token)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf(tokenString)

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	//} else {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//	m, _ := json.Marshal("Either the Username or Password was invalid")
	//	_, _ = w.Write(m)
	////}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code up until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	fmt.Printf("1")
	if err != nil {
		fmt.Printf("2")
		if err == http.ErrNoCookie {
			fmt.Printf("no cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Printf("3")
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		fmt.Printf("4")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code up until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf(tokenString)

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &models.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
/**************************************************************/



func main() {
	// "Signin" and "Welcome" are the actions that we will implement
	router := mux.NewRouter()
	router.HandleFunc("/api/signin", Signin)
	router.HandleFunc("/api/welcome", Welcome)
	router.HandleFunc("/api/refresh", Refresh)
	router.HandleFunc("/api/all_students", StudentsIndex)
	router.HandleFunc("/api/students", CreateStudent).Methods(http.MethodPost)
	router.HandleFunc("/api/students/{studentId}", FetchStudent).Methods(http.MethodGet)
	//router.HandleFunc("/api/students/{studentId}", UpdateStudent).Methods(http.MethodPut)
	router.HandleFunc("/api/students/{studentId}", DeleteStudent).Methods(http.MethodDelete)

	// start the server on port 8000

	log.Fatal(http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
}