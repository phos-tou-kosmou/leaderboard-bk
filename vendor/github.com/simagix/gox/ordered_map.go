// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import (
	"bytes"
	"encoding/json"
	"sort"
)

// OrderedMap preserves keys order
type OrderedMap struct {
	SortedKeys []string
	Map        map[string]interface{}
}

// NewOrderedMap returns an ordered map
func NewOrderedMap(str string) *OrderedMap {
	var om OrderedMap
	json.Unmarshal([]byte(str), &om)
	return &om
}

// UnmarshalJSON is used by json.Unmarshal
func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	json.Unmarshal(b, &om.Map)
	index := make(map[string]int)
	for key := range om.Map {
		om.SortedKeys = append(om.SortedKeys, key)
		esc, _ := json.Marshal(key)
		index[key] = bytes.Index(b, esc)
	}
	sort.Slice(om.SortedKeys, func(i, j int) bool { return index[om.SortedKeys[i]] < index[om.SortedKeys[j]] })
	return nil
}

// MarshalJSON is used by json.Marshal
func (om OrderedMap) MarshalJSON() ([]byte, error) {
	var err error
	var b, kb, vb []byte
	buffer := bytes.NewBuffer(b)
	buffer.WriteRune('{')
	for i, key := range om.SortedKeys {
		if kb, err = json.Marshal(key); err != nil {
			return nil, err
		}
		buffer.Write(kb)
		buffer.WriteRune(':')
		if vb, err = json.Marshal(om.Map[key]); err != nil {
			return nil, err
		}
		buffer.Write(vb)
		if i != len(om.SortedKeys)-1 {
			buffer.WriteRune(',')
		}
	}
	buffer.WriteRune('}')
	return buffer.Bytes(), nil
}
