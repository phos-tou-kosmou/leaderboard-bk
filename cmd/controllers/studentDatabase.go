package controllers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	_ "sync"
	"time"
)

type Student struct {
	ID int
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	GPA float32 `json:"gpa"`
	Sport string `json:"sport"`
	CreatedAt time.Time
}

func dbConn() (db *sql.DB) {
	//dbDriver := "mysql"
	//dbUser := "root"
	//dbPass := "root"
	//dbName := "leaderboard"
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/leaderboard")
	if err != nil {
		log.Panic(err.Error())
	}
	return db
}

/******************************************************************************/

func IndexStudents(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM leaderboard.students")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	stus := make([]*Student, 0)
	for rows.Next() {
		stu := new(Student)
		err := rows.Scan(&stu.ID, &stu.FirstName, &stu.LastName, &stu.GPA, &stu.Sport, &stu.CreatedAt)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		stus = append(stus, stu)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, stu := range stus {
		result, err := fmt.Fprint(w, "%d, %s, %s, %.2f, %s", stu.ID, stu.FirstName, stu.LastName, stu.GPA, stu.Sport )
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result)
		_, _ = w.Write([]byte(`{"test":"this is a test"}`))
	}
}

/*******************************************************************************/

func EditStudent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM leaderboard.students", nId)
	if err != nil {
		panic(err.Error())
	}
	stu := Student{}
	for selDB.Next() {
		var firstName, lastName, sport string
		var gpa float32
		err = selDB.Scan(&firstName, &lastName, &gpa, &sport)
		if err != nil {
			panic(err.Error())
		}
		stu.FirstName = firstName
		stu.LastName = lastName
		stu.GPA = gpa
		stu.Sport = sport
	}
	defer db.Close()
}

/******************************************************************************/

func InsertStudent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("I am here")
	if r.Method == "POST" {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		gpa := r.FormValue("gpa")
		sport := r.FormValue("sport")
		insForm, err := db.Prepare("INSERT INTO leaderboard.students(first_name, last_name, gpa, sport) VALUES(?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		_, _ = insForm.Exec(firstName, lastName, gpa, sport)
		log.Println(
			"INSERT: First Name: " + firstName +
			" | Last Name: " + lastName +
			" | GPA: " + gpa +
			" | Sport " + sport)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/******************************************************************************/

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Now I am here")
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE leaderboard.students SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		_, _ = insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/******************************************************************************/

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	stu := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM leaderboard.students WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	_, _ = delForm.Exec(stu)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/******************************************************************************/

