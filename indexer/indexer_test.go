package indexer

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	i := Init()
	id, _ := GeneDocId(12, 1, 11)
	i.addDoc(&Doc{
		DocId: id,
		TimeStamp:123,
	}, "测试", false)
	id1, _ := GeneDocId(11, 1, 11)
	i.addDoc(&Doc{
		DocId:id1,
		TimeStamp: 123444,
	}, "测试", false)
	i.FlushDisk()
	fmt.Println("测试")
}

func TestIndexer_LoadFromDisk(t *testing.T) {
	fmt.Println(uint16(0))
	i := Init()
	i.LoadFromDisk()
	fmt.Println("iii")
}

func TestGeneDocId(t *testing.T) {
	s := 10
	e := 89

	m := (s + e) / 2
	fmt.Println("m", m)
}
