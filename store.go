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
	PERSIST_MANUAL = iota
	PERSIST_WRITES
	PERSIST_INTERVAL
)

var (
	ErrInvalidPersist = errors.New("invalid persist type")
)

type store struct {
	sync.RWMutex

	filename string
	items    []Item
	persist  int
}

type Type struct {
	Name  string
	Value interface{}
}

func NewStore(filename string, persist int, types ...Type) (*store, error) {
	for _, t := range types {
		gob.RegisterName(t.Name, t.Value)
	}

	if persist != PERSIST_MANUAL &&
		persist != PERSIST_WRITES &&
		persist != PERSIST_INTERVAL {
		return nil, ErrInvalidPersist
	}

	store := &store{
		filename: filename,
		items:    make([]Item, 0),
		persist:  persist,
	}

	if persist == PERSIST_INTERVAL {
		go store.persistInterval()
	}

	return store, nil
}

type Item interface {
	MarshalJSON() ([]byte, error)
}

func (s *store) persistInterval() {
	for {
		time.Sleep(time.Second * 60)

		println("writing")

		s.Lock()
		if err := s.Save(); err != nil {
			s.Unlock()
			panic(err)
		}
		s.Unlock()
	}
}

func (s *store) Add(item Item) (int, error) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, item)

	if s.persist == PERSIST_WRITES {
		return len(s.items) - 1, s.Save()
	}

	return len(s.items) - 1, nil
}

func (s *store) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`[`)
	for i, item := range s.items {
		d, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buffer.Write(d)
		if i != len(s.items)-1 {
			buffer.WriteString(`,`)
		}
	}
	buffer.WriteString(`]`)

	return buffer.Bytes(), nil
}

func (s *store) encodeGob() []byte {
	var buffer bytes.Buffer

	gob.NewEncoder(&buffer).Encode(s.items)

	return buffer.Bytes()
}

func (s *store) decodeGob(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&s.items)
}

func (s *store) Save() error {
	return ioutil.WriteFile(s.filename, s.encodeGob(), 0644)
}

func (s *store) Load() error {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(s.filename); err != nil {
		return err
	}

	return s.decodeGob(data)
}

func (s *store) Update(f func(items []Item) error) error {
	s.Lock()
	defer s.Unlock()

	if err := f(s.items); err != nil {
		return err
	}

	if s.persist == PERSIST_WRITES {
		return s.Save()
	}

	return nil
}

func (s *store) View(f func(items []Item)) {
	s.RLock()
	defer s.RUnlock()

	f(s.items)
}
