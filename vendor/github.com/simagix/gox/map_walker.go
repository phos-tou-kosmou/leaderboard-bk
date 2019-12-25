// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import (
	"reflect"
)

type callback func(interface{}) interface{}

// MapWalker is an empty JSON document
type MapWalker struct {
	cb callback
}

// NewMapWalker returns a MapWalker
func NewMapWalker(cb callback) *MapWalker {
	return &MapWalker{cb: cb}
}

// Walk walks a map
func (walker *MapWalker) Walk(docMap map[string]interface{}) map[string]interface{} {
	for k, v := range docMap {
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if mv, ok := v.(map[string]interface{}); ok {
				docMap[k] = walker.Walk(mv)
			} else {
				panic(v)
			}
		case reflect.Array, reflect.Slice:
			if mv, ok := v.([]interface{}); ok {
				docMap[k] = walker.iterate(mv)
			} else {
				panic(v)
			}
		default:
			if walker.cb != nil {
				docMap[k] = walker.cb(v)
			}
		}
	}
	return docMap
}

// iterate iterates thru an array
func (walker *MapWalker) iterate(arrayType []interface{}) []interface{} {
	for k, v := range arrayType {
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if mv, ok := v.(map[string]interface{}); ok {
				arrayType[k] = walker.Walk(mv)
			} else {
				panic(v)
			}
		case reflect.Array, reflect.Slice:
			if mv, ok := v.([]interface{}); ok {
				arrayType[k] = walker.iterate(mv)
			} else {
				panic(v)
			}
		default:
			arrayType[k] = walker.cb(v)
		}

	}
	return arrayType
}
