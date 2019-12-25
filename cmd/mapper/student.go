package mapper

import (
	"encoding/json"
	"leaderboard-bk/cmd/models"
	"net/http"
)

var stu models.Student
type StudentMap struct { }

func (stu *StudentMap) Insert(p *models.Student) error {
	result, err := db.Exec("INSERT INTO books VALUES($1, $2, $3, $4)", isbn, title, author, price)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return err
	}

	var w http.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	m, _ := json.Marshal(result)
	_, _ = w.Write(m)
	return nil
}

func (stu *StudentMap) Update(p *models.Student) error {


	return nil
}

func (stu *StudentMap) Delete(p *models.Student) error {


	return nil
}