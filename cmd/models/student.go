package models

import (
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

type Golfer struct {
	// Pointer to student struct to mimic inheritance
	// TODO: Need to look at Context and Reflect
	*Student
	ShortCourseAvg float32
	FullCourseAvg float32
}
