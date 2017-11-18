// Package store manages a list of items in memory with concurrency-safe
// functions for accessing and manipulating them.
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
	// PERSIST_MANUAL is the default disk-persistence setting. The user has to
	// manually save the store to disk using the Save() function.
	PERSIST_MANUAL = iota

	// PERSIST_WRITES is a disk-persistence setting that will save the store
	// to disk when either Add() or Update() is called.
	PERSIST_WRITES

	// PERSIST_INTERVAL is a disk-persistence setting that will save the store
	// on a set interval (currently once every 60 seconds).
	PERSIST_INTERVAL
)

var (
	// ErrInvalidPersist is thrown when an invalid disk-persistence setting is
	// provided when calling NewStore().
	ErrInvalidPersist = errors.New("invalid persist type")

	// ErrInvalidTypes is thrown when no types are listed when calling
	// NewStore(). Without type definitions it is impossible to persist on disk
	// without possible data loss.
	ErrInvalidTypes = errors.New("no types were provided")

	// ErrInvalidFilename is thrown when NewStore() is called with an empty
	// filename - making persistence impossible.
	ErrInvalidFilename = errors.New("invalid filename")
)

// Store represents a collection of items that persist on disk.
type Store struct {
	filename string
	items    []Item
	persist  int
	mutex    sync.RWMutex
}

// Type is used to register types from outside packages so that they are
// recognized when loading or saving the store.
type Type struct {
	// Name is the name of this type. Usually it takes the form of
	// "package.Name" -- so for a struct User{} in package main this field
	// would be "main.User".
	Name string

	// Value is an empty struct of this type. For a struct User{}, this would
	// be User{}.
	Value interface{}
}

// NewStore is the primary constructor function for creating stores. The
// provided filename is where the store will persist to disk (and read from
// disk). The persist int is one of the store.PERSIST_ constants. The provided
// types register the types that will be held in the store.
//
// NewStore will return an error if the persist parameter is not a valid
// store.PERSIST_ constant.
func NewStore(filename string, persist int, types ...Type) (*Store, error) {
	if len(filename) == 0 {
		return nil, ErrInvalidFilename
	}

	if len(types) == 0 {
		return nil, ErrInvalidTypes
	}

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

// Item implements the json.Marshaler interface and is used so that the store
// itself can implement the json.Marshaler function by aggregating all items.
type Item interface {
	MarshalJSON() ([]byte, error)
}

func (s *Store) persistInterval() {
	for {
		time.Sleep(time.Second * 60)

		if err := s.Save(); err != nil {
			println(err.Error())
		}
	}
}

// Add appends an Item on the end of the store. It returns the id of the item
// and an error if there was a problem persisting the store on the disk (if
// PERSIST_WRITE is enabled).
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

// MarshalJSON returns the store as a JSON list. It returns an error if there
// was an error marshaling one of the items.
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

// Save persists the store on disk using the filename provided when NewStore()
// was called.
func (s *Store) Save() error {
	return ioutil.WriteFile(s.filename, s.encodeGob(), 0644)
}

// Load reads the store from disk using the filename provided when NewStore()
// was called.
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

// Update is used to manipulate an item (or items) in the store. It returns
// an error if there is an error saving the store (if PERSIST_WRITES is
// enabled) or if there is an error inside the f function.
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

// View is used to read an item (or items) in the store.
func (s *Store) View(f func(items []Item)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	f(s.items)
}
