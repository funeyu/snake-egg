package store

import (
	"fmt"
	"testing"
)

func TestBadger_Add(t *testing.T) {
	b := InitBadger("./badger")

	g := b.Get("18439249802455220228") // 18437349846362030087
	fmt.Println("g", g)
}
