// Copyright 2019 Kuei-chun Chen. All rights reserved.

package gox

import "encoding/json"

// Stringify returns a string of a map
func Stringify(doc interface{}, opts ...string) string {
	var err error
	var data []byte
	if doc == nil {
		return ""
	} else if len(opts) == 2 {
		if data, err = json.MarshalIndent(doc, opts[0], opts[1]); err != nil {
			return ""
		}
		return string(data)
	} else if data, err = json.Marshal(doc); err != nil {
		return ""
	}
	return string(data)
}
