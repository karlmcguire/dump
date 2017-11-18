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

type Store struct {
	filename string
	items    []Item
	persist  int
	mutex    sync.RWMutex
}

type Type struct {
	Name  string
	Value interface{}
}

func NewStore(filename string, persist int, types ...Type) (*Store, error) {
	for _, t := range types {
		gob.RegisterName(t.Name, t.Value)
	}

	if persist != PERSIST_MANUAL &&
		persist != PERSIST_WRITES &&
		persist != PERSIST_INTERVAL {
		return nil, ErrInvalidPersist
	}

	store := &Store{
		filename: filename,
		items:    make([]Item, 0),
		persist:  persist,
		mutex:    sync.RWMutex{},
	}

	if persist == PERSIST_INTERVAL {
		go store.persistInterval()
	}

	return store, nil
}

type Item interface {
	MarshalJSON() ([]byte, error)
}

func (s *Store) persistInterval() {
	for {
		time.Sleep(time.Second * 60)

		println("writing")

		s.mutex.Lock()
		if err := s.Save(); err != nil {
			s.mutex.Unlock()
			panic(err)
		}
		s.mutex.Unlock()
	}
}

func (s *Store) Add(item Item) (int, error) {
	s.mutex.Lock()

	s.items = append(s.items, item)

	if s.persist == PERSIST_WRITES {
		s.mutex.Unlock()
		return len(s.items) - 1, s.Save()
	}

	s.mutex.Unlock()
	return len(s.items) - 1, nil
}

func (s *Store) MarshalJSON() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

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

func (s *Store) encodeGob() []byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var buffer bytes.Buffer

	gob.NewEncoder(&buffer).Encode(s.items)

	return buffer.Bytes()
}

func (s *Store) decodeGob(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&s.items)
}

func (s *Store) Save() error {
	return ioutil.WriteFile(s.filename, s.encodeGob(), 0644)
}

func (s *Store) Load() error {
	s.mutex.Lock()

	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(s.filename); err != nil {
		s.mutex.Unlock()
		return err
	}

	s.mutex.Unlock()
	return s.decodeGob(data)
}

func (s *Store) Update(f func(items []Item) error) error {
	s.mutex.Lock()

	if err := f(s.items); err != nil {
		s.mutex.Unlock()
		return err
	}

	if s.persist == PERSIST_WRITES {
		s.mutex.Unlock()
		return s.Save()
	}

	s.mutex.Unlock()
	return nil
}

func (s *Store) View(f func(items []Item)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	f(s.items)
}
