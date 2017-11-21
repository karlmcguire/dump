package main

import (
	"bytes"
)

// Post is just a sample struct.
type Post struct {
	Name string
	Body string
}

// Post implements the json.Marshaler interface.
func (p *Post) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"name":"`)
	buffer.WriteString(p.Name)
	buffer.WriteString(`","body":"`)
	buffer.WriteString(p.Body)
	buffer.WriteString(`"}`)

	return buffer.Bytes(), nil
}
