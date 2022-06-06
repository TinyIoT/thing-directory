package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDB storage
type LevelDBStorage struct {
	db *leveldb.DB
	wg sync.WaitGroup
}

func NewLevelDBStorage(dsn string, opts *opt.Options) (Storage, error) {
	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	// Open the database file
	db, err := leveldb.OpenFile(url.Path, opts)
	if err != nil {
		return nil, err
	}

	return &LevelDBStorage{db: db}, nil
}

// CRUD
func (s *LevelDBStorage) add(id string, td ThingDescription) error {
	if id == "" {
		return fmt.Errorf("ID is not set")
	}

	bytes, err := json.Marshal(td)
	if err != nil {
		return err
	}

	found, err := s.db.Has([]byte(id), nil)
	if err != nil {
		return err
	}
	if found {
		return &ConflictError{id + " is not unique"}
	}

	err = s.db.Put([]byte(id), bytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *LevelDBStorage) get(id string) (ThingDescription, error) {

	bytes, err := s.db.Get([]byte(id), nil)
	if err == leveldb.ErrNotFound {
		return nil, &NotFoundError{id + " is not found"}
	} else if err != nil {
		return nil, err
	}

	var td ThingDescription
	err = json.Unmarshal(bytes, &td)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *LevelDBStorage) update(id string, td ThingDescription) error {

	bytes, err := json.Marshal(td)
	if err != nil {
		return err
	}

	found, err := s.db.Has([]byte(id), nil)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{id + " is not found"}
	}

	err = s.db.Put([]byte(id), bytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *LevelDBStorage) delete(id string) error {
	found, err := s.db.Has([]byte(id), nil)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{id + " is not found"}
	}

	err = s.db.Delete([]byte(id), nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *LevelDBStorage) listPaginate(offset, limit int) ([]ThingDescription, error) {

	// TODO: is there a better way to do this?
	TDs := make([]ThingDescription, 0, limit)
	s.wg.Add(1)
	iter := s.db.NewIterator(nil, nil)

	for i := 0; i < offset+limit && iter.Next(); i++ {
		if i >= offset && i < offset+limit {
			var td ThingDescription
			err := json.Unmarshal(iter.Value(), &td)
			if err != nil {
				return nil, err
			}
			TDs = append(TDs, td)
		}
	}
	iter.Release()
	s.wg.Done()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return TDs, nil
}

func (s *LevelDBStorage) listAllBytes() ([]byte, error) {

	s.wg.Add(1)
	iter := s.db.NewIterator(nil, nil)

	var buffer bytes.Buffer
	buffer.WriteString("[")
	separator := byte(',')
	first := true
	for iter.Next() {
		if first {
			first = false
		} else {
			buffer.WriteByte(separator)
		}
		buffer.Write(iter.Value())
	}
	buffer.WriteString("]")

	iter.Release()
	s.wg.Done()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *LevelDBStorage) iterate() <-chan ThingDescription {
	serviceIter := make(chan ThingDescription)

	go func() {
		defer close(serviceIter)

		s.wg.Add(1)
		defer s.wg.Done()
		iter := s.db.NewIterator(nil, nil)
		defer iter.Release()

		for iter.Next() {
			var td ThingDescription
			err := json.Unmarshal(iter.Value(), &td)
			if err != nil {
				log.Printf("LevelDB Error: %s", err)
				return
			}
			serviceIter <- td
		}

		err := iter.Error()
		if err != nil {
			log.Printf("LevelDB Error: %s", err)
		}
	}()

	return serviceIter
}

func (s *LevelDBStorage) iterateBytes(ctx context.Context) <-chan []byte {
	bytesCh := make(chan []byte, 0) // must be zero

	go func() {
		defer close(bytesCh)

		s.wg.Add(1)
		defer s.wg.Done()
		iter := s.db.NewIterator(nil, nil)
		defer iter.Release()

	Loop:
		for iter.Next() {
			select {
			case <-ctx.Done():
				//log.Println("LevelDB: canceled")
				break Loop
			default:
				b := make([]byte, len(iter.Value()))
				copy(b, iter.Value())
				bytesCh <- b
			}
		}

		err := iter.Error()
		if err != nil {
			log.Printf("LevelDB Error: %s", err)
		}
	}()

	return bytesCh
}

func (s *LevelDBStorage) Close() {
	s.wg.Wait()
	err := s.db.Close()
	if err != nil {
		log.Printf("Error closing storage: %s", err)
	}
	if flag.Lookup("test.v") == nil {
		log.Println("Closed leveldb.")
	}
}
