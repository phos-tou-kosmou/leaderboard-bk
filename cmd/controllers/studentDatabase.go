package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"leaderboard-bk/cmd/models"
	"log"
	"net/http"
	_ "sync"
	_ "time"
)



func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql",
		"root:root@tcp(localhost:3306)/leaderboard")
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

	stus := make([]*models.Student, 0)
	var newTime *mysql.NullTime
	for rows.Next() {
		stu := new(models.Student)
		err := rows.Scan(&stu.ID,
			&stu.FirstName,
			&stu.LastName,
			&stu.GPA,
			&stu.Sport,
			&newTime)
		log.Println(err)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		stus = append(stus, stu)
		//log.Println(stu.FirstName, stu.LastName, stu.GPA, stu.Sport, stu.CreatedAt)
	}
	//if err = rows.Err(); err != nil {
	//	http.Error(w, http.StatusText(500), 500)
	//	return
	//}

	//for _, stu := range stus {
	//	result, err := fmt.Fprint(w, "%d, %s, %s, %d, %s", stu.ID, stu.FirstName, stu.LastName, stu.GPA, stu.Sport)
	//	if err != nil {
	//		http.Error(w, http.StatusText(500), 500)
	//		log.Println(result)
	//		return
	//	}
	//
	//}

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	//w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"hello":"world"}`))
}

/*******************************************************************************/

func ModifyStudent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM leaderboard.students", nId)
	if err != nil {
		panic(err.Error())
	}
	stu := models.Student{}
	for selDB.Next() {
		var ID int
		var firstName, lastName, sport string
		var gpa float32
		err = selDB.Scan(&ID, &firstName, &lastName, &gpa, &sport)
		if err != nil {
			panic(err.Error())
		}
		stu.ID = ID
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
	if r.Method == "POST" {
		var s *models.Student_test
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(s.FirstName)
		firstName := s.FirstName
		lastName := s.LastName
		gpa := s.GPA
		sport := s.Sport
		insForm, err :=
			db.Prepare("INSERT INTO leaderboard.students(firstName, lastName, gpa, sport) VALUES(?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		_, _ = insForm.Exec(firstName, lastName, gpa, sport)
		fts := fmt.Sprintf("f%",  gpa)
		log.Println(
			"INSERT: First Name: " + firstName +
			" | Last Name: " + lastName +
			" | GPA: " + fts +
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
		insForm, err := db.Prepare("UPDATE leaderboard.students SET first_name=?, last_name=?, gpa=?, sport=?,  WHERE not last_name=?")
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

