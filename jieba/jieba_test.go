package jieba

import (
	"fmt"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	s := " )("
	ss := Cut(s)
	fmt.Println("s", ss)
	var ls = "���"
	fmt.Println(strings.Contains(ls, "�"))
}

func TestSplitWords(t *testing.T) {
}