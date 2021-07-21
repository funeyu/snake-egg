package jieba

import (
	"fmt"
	"snake/category"
	"snake/db/models"
	"snake/util"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	s := "小人物 - Nobody 60FPS/帧率电影高速下载"
	ss := Cut(s)
	fmt.Println("s", ss)
}

func TestSplitWords(t *testing.T) {
	sites := models.AllSites()
	words := make(map[string] []string)
	for _, site := range sites {
		if strings.Contains(site.Keywords, "�") {
			continue
		}
		if !util.IsChinese(site.Keywords) {
			_, ok := words["外文网站"]
			if !ok {
				words["外文网站"] = []string{site.Url}
			} else {
				words["外文网站"] = append(words["外文网站"], site.Url)
			}
			continue
		}

		ss := Cut(site.Keywords)
		cate := category.Sort(ss)
		if cate != "" {
			_, ok := words[cate]
			if ok {
				words[cate] = append(words[cate], site.Url)
			} else {
				words[cate] = []string{site.Url}
			}
		}
	}
	fmt.Println(words)
	}