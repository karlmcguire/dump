package main

import (
	"bytes"
)

type Post struct {
	Name string
	Body string
}

func (p *Post) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"name":"`)
	buffer.WriteString(p.Name)
	buffer.WriteString(`","body":"`)
	buffer.WriteString(p.Body)
	buffer.WriteString(`"}`)

	return buffer.Bytes(), nil
}
