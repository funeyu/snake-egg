package jieba

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type dict struct {
	words     map[string]int64
	totalFreq int64
}

var min_FREP = -1e+100 //最小的频次

func load(filename string) *dict {
	d := &dict{
		words:     make(map[string]int64, 1024),
		totalFreq: 0,
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), " ")[:2]
		i, e := strconv.ParseInt(splits[1], 10, 0)
		if e == nil {
			d.totalFreq = d.totalFreq + i
			d.words[splits[0]] = i
		}
	}
	return d
}

func (d *dict) contains(word string) bool {
	_, ok := d.words[word]
	return ok
}

// 频率等于该词在语料库中出现的频次 / 总频次，为了方便计算 取了log2
func (d *dict) freq(word string) float64 {
	freq, ok := d.words[word]
	if !ok {
		return min_FREP
	}

	return math.Log2(float64(freq)) - math.Log2(float64(d.totalFreq))
}
