package store

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"sync"
)

type Store struct {
	sync.RWMutex

	filename string
	items    []Item
}

func NewStore(filename string) *Store {
	return &Store{
		filename: filename,
		items:    make([]Item, 0),
	}
}

type Item interface {
	EncodeJson() []byte
}

func (s *Store) Add(item Item) int {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)

	return len(s.items) - 1
}

func (s *Store) EncodeJson() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for i, item := range s.items {
		buffer.Write(item.EncodeJson())
		if i != len(s.items)-1 {
			buffer.WriteString(`,`)
		}
	}
	buffer.WriteString(`]`)

	return buffer.Bytes()
}

func (s *Store) encodeGob() []byte {
	var buffer bytes.Buffer

	gob.NewEncoder(&buffer).Encode(s.items)

	return buffer.Bytes()
}

func (s *Store) decodeGob(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&s.items)
}

func (s *Store) EncodeFile() error {
	return ioutil.WriteFile(s.filename, s.encodeGob(), 0644)
}

func (s *Store) DecodeFile() error {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(s.filename); err != nil {
		return err
	}

	return s.decodeGob(data)
}

func (s *Store) Write(f func(items []Item) error) error {
	s.Lock()
	defer s.Unlock()

	return f(s.items)
}

func (s *Store) Read(f func(items []Item)) {
	s.RLock()
	defer s.RUnlock()

	f(s.items)
}
