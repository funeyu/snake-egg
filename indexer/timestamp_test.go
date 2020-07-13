package indexer

import (
	"fmt"
	"testing"
)

func TestTimestampIndex_IterateAll(t *testing.T) {
	time := GenerateTimestampIndex()
	id,_ := GeneDocId(1, 12, 12)
	time.Add(id, 12345678)

	id1, _ := GeneDocId(1, 12, 13)
	time.Add(id1, 234)

	id2, _:= GeneDocId(1, 12, 14)
	time.Add(id2, 3456788745)
	for _, t := range time.IterateAll() {
		fmt.Println("t", t)
	}
}
