package store

import (
	"bytes"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	type params struct {
		filename string
		persist  int
		types    []Type
		err      error
	}

	tests := []params{
		params{
			filename: "test.db",
			persist:  PERSIST_WRITES,
			types:    nil,
			err:      ErrInvalidTypes,
		},
		params{
			filename: "test.db",
			persist:  999,
			types:    []Type{Type{"meh", params{}}},
			err:      ErrInvalidPersist,
		},
		params{
			filename: "test.db",
			persist:  PERSIST_WRITES,
			types:    []Type{Type{"meh", params{}}},
			err:      nil,
		},
		params{
			filename: "",
			persist:  PERSIST_WRITES,
			types:    []Type{Type{"meh", params{}}},
			err:      ErrInvalidFilename,
		},
		params{
			filename: "test.db",
			persist:  PERSIST_INTERVAL,
			types:    []Type{Type{"meh", params{}}},
			err:      nil,
		},
	}

	var err error

	for _, v := range tests {
		if v.types != nil {
			if _, err = NewStore(v.filename, v.persist, v.types[0]); err != v.err {
				t.Fatal("unexpected error")
			}
		} else {
			if _, err = NewStore(v.filename, v.persist); err != v.err {
				t.Fatal("unexpected error")
			}
		}
	}
}

type Blob struct {
	Data string
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"data":"`)
	buffer.WriteString(b.Data)
	buffer.WriteString(`"}`)

	return buffer.Bytes(), nil
}

func TestPersistInterval(t *testing.T) {
	test, _ := NewStore("test.db", PERSIST_INTERVAL, Type{"store.Blob", &Blob{}})

	_, _ = test.Add(&Blob{"meh"})

	time.Sleep(time.Second * 70)
}
