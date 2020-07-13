package store

import "snake/indexer"

func Init() *Memory {
	return &Memory{
		Docs:make(map[int64]indexer.Doc),
	}
}

type Memory struct {
	Docs map[int64]indexer.Doc
}

func(m *Memory) Get(docId int64) indexer.Doc {
	return m.Docs[docId]
}

func (m *Memory) Add(doc indexer.Doc) {
	m.Docs[int64(doc.DocId)] = doc
}

func (m *Memory) Prepare() error {
	return nil
}

func (m *Memory) Commit() error {
	return nil
}
