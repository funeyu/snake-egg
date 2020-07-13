package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"log"
	"snake/indexer"
)

func InitBadger(path string) *Badger{
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal("initBadger with fatal error:", err)
	}

	return &Badger{
		path: path,
		docStore:db,
	}
}

type Badger struct {
	path string
	docStore *badger.DB
}
func encode(doc indexer.Doc) []byte {
	buffer := new (bytes.Buffer)
	error := json.NewEncoder(buffer).Encode(doc)
	if error != nil {
		return nil
	}
	return buffer.Bytes()
}
func (b *Badger) Add(doc indexer.Doc) error {
	txn := b.docStore.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set([]byte(doc.Id), encode(doc))
	if err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func decode(bytes []byte) *indexer.Doc {
	d := indexer.Doc {}
	err := json.Unmarshal(bytes, &d)
	if err != nil {
		fmt.Println("decode to doc error:", err)
		return nil
	}
	return &d
}

func (b *Badger) Get(docId string) indexer.Doc {
	var valCopy []byte
	b.docStore.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte(docId))
		if item == nil {
			return nil
		}
		item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return nil
	})
	if valCopy == nil {
		return indexer.Doc{}
	}
	return *decode(valCopy)
}

func (b *Badger) Update(doc indexer.Doc) error {
	d := b.Get(doc.Id)
	if d == (indexer.Doc{}) {
		return errors.New("can not find docId:" + doc.Id)
	}

	b.Add(doc)
	return nil
}

func (b *Badger) Commit() error {
	err := b.docStore.Close()
	return err
}

func (b *Badger) ForEach(fn func(docId string, v string)error) error {
	return nil
}

func (b *Badger) Prepare() error {
	return nil
}
