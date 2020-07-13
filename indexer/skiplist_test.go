package indexer

import (
	"encoding/binary"
	"fmt"
	"testing"
)
type IntElement struct {
	i uint
}

func (i IntElement) Compare(other interface{}) int {
	return int(i.i)  - int(other.(IntElement).i)
}

func (i IntElement) Serialize() []byte {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(i.i))
	return out
}

func (i IntElement) ByteSize() uint8 {
	return 4
}

func (i IntElement) UnSerialize(bytes []byte) {
	i.i = uint(binary.BigEndian.Uint32(bytes))
}

func TestIsEmpty(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})
	fmt.Println("len(15)", s.isEmpty())
}

func TestSkipList_Insert(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})
	s.Insert(IntElement{i: 1})
	s.Insert(IntElement{i: 5})
	s.Insert(IntElement{i: 8})
	fmt.Println(s)
}

func TestSkipList_Find(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})
	for i := 0; i < 7000; i ++ {
		s.Insert(IntElement{i: uint(i)})
	}
	s.Delete(IntElement{i: uint(10)})

	for i := 0; i < 700; i ++ {
		found := s.Find(IntElement{i: uint(i)})
		if found != nil {
			fmt.Println("found", found.value)
		}
	}
}

func TestLevel(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})
	for i:=0; i < 1000; i ++ {
		fmt.Println(s.level())
	}
}

func TestSkipList_Iterator(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})


	s.Insert(IntElement{234})
	s.Insert(IntElement{34})
	s.Insert(IntElement{2333})


	fmt.Println("iterator", s.Iterator(3))
}

func TestSkipList_SerializeToBytes(t *testing.T) {
	s := Generate(15, func(bytes []byte) Element {
		return nil
	})

	for i := 0; i < 700000; i ++ {
		s.Insert(IntElement{i: uint(i)})
	}

	bytes := s.SerializeToBytes()
	skip := SkipListFromBytes(bytes, func(bytes []byte) Element {
		i := IntElement{
			i:uint(binary.BigEndian.Uint32(bytes)),
		}
		return i
	}, 4)

	fmt.Println("77", skip.Find(IntElement{i: uint(6000)}))
}
