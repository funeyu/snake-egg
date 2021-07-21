package store

import (
	"snake/indexer")

// 先最简单的以string来实现
type  Store interface {
	Add(doc indexer.Doc) error
	Get(docId string) indexer.Doc
	Prepare() error
	Commit() error
	Update(doc indexer.Doc) error
	ForEach(fn func(id int, v string)error) error
}
