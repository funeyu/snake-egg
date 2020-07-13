package resource

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func DictWords() map[string]bool{
	words := make(map[string]bool)

	wd := "/Users/fuheyu/egg-crawler"
	file, err := os.Open(wd + "/jieba/resource/words.data")
	if err != nil {
		panic("err when open file ")
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	re := regexp.MustCompile(`[\s]+`)

	for fileScanner.Scan() {
		t := fileScanner.Text()
		ss := re.Split(t, -1)
		word := ss[0]
		word = strings.Replace(word, "的", "", -1)
		words[word] =  true
	}

	return words
}