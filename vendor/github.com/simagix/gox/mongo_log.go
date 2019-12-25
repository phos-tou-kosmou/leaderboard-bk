// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import "strings"

// MongoLog stores a line of mongo log
type MongoLog struct {
	log string
}

// NewMongoLog returns a mongo log struct
func NewMongoLog(log string) *MongoLog {
	return &MongoLog{log: log}
}

// Get returns a JSON string of a tag
func (ml *MongoLog) Get(key string) string {
	str := ml.log
	i := strings.Index(str, key)
	if i < 0 {
		return ""
	}
	str = strings.Trim(str[i+len(key):], " ")
	isFound := false
	bpos := 0 // begin position
	epos := 0 // end position
	for _, r := range str {
		epos++
		if isFound == false && r == '{' {
			isFound = true
			bpos++
		} else if isFound == true {
			if r == '{' {
				bpos++
			} else if r == '}' {
				bpos--
			}
		}

		if isFound == true && bpos == 0 {
			break
		}
	}
	return str[bpos:epos]
}
