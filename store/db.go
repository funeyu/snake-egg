package store

import (
	"snake/db/models"
	"snake/indexer"
)

type DBStore struct {
	details []models.CrawlDetail
}

func InitDBStore() *DBStore {
	return &DBStore{
		details: models.AllCrawlDetails(),
	}
}

func (ds *DBStore) Get(docId string) indexer.Doc {
	return indexer.Doc{}
}

func (ds *DBStore) Add(doc indexer.Doc) error { return nil }

func (ds *DBStore) Prepare() error { return nil}

func (ds *DBStore) Commit() error {return nil}

func (ds *DBStore) Update(doc indexer.Doc) error {return nil}

func (ds *DBStore) ForEach(fn func(id int, v string)error) error {
	for _, detail := range ds.details {
		fn(detail.ID, detail.Keyword)
	}
	return nil
}


