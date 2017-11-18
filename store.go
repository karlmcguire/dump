package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

const (
	MANUAL = iota
	WRITES
	INTERVAL
)

var (
	ErrInvalidPersist = errors.New("invalid persist type")
)

type Store struct {
	sync.RWMutex

	filename string
	items    []Item
	persist  int
}

type Type struct {
	Name  string
	Value interface{}
}

func NewStore(filename string, persist int, types ...Type) (*Store, error) {
	for _, t := range types {
		gob.RegisterName(t.Name, t.Value)
	}

	if persist != MANUAL && persist != WRITES && persist != INTERVAL {
		return nil, ErrInvalidPersist
	}

	store := &Store{
		filename: filename,
		items:    make([]Item, 0),
		persist:  persist,
	}

	if persist == INTERVAL {
		go store.persistInterval()
	}

	return store, nil
}

type Item interface {
	EncodeJson() []byte
}

func (s *Store) persistInterval() {
	for {
		time.Sleep(time.Second * 60)

		println("writing")

		s.Lock()
		if err := s.EncodeFile(); err != nil {
			s.Unlock()
			panic(err)
		}
		s.Unlock()
	}
}

func (s *Store) Add(item Item) (int, error) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)

	if s.persist == WRITES {
		return len(s.items) - 1, s.EncodeFile()
	}

	return len(s.items) - 1, nil
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

	if err := f(s.items); err != nil {
		return err
	}

	if s.persist == WRITES {
		return s.EncodeFile()
	}

	return nil
}

func (s *Store) Read(f func(items []Item)) {
	s.RLock()
	defer s.RUnlock()

	f(s.items)
}
