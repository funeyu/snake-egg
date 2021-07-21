package models

import (
	"encoding/json"
	"snake/db"
	"strings"
)

type Keywords struct {
	ID int `gorm:"primary_key"`
	Key string `gorm:"type:varchar(16);"`
	Words string `gorm:"type:varchar(256);"`
}

type Word struct {
	Logic string `json:"logic"`
	Values []string `json:"values"`
}

type KeywordInfo struct {
	Keyword string
	Words []Word
}

func allKeywords() []KeywordInfo{
	var keys []Keywords
	db.DB.Find(&keys)

	var res []KeywordInfo
	for _, kw := range keys {
		var words []Word
		if kw.Words != "" {
			json.Unmarshal([]byte(kw.Words), &words)
		}
		res = append(res, KeywordInfo{
			Keyword: kw.Key,
			Words: words,
		})
	}
	return res
}

var keys []KeywordInfo

func init() {
	keys = allKeywords()
}

func orLogic(s string, w Word) bool {
	for _, v := range w.Values {
		if strings.Contains(s, v) {
			return true
		}
	}

	return false
}

func andLogin(s string, w Word) bool {
	for _, v := range w.Values {
		if !strings.Contains(s, v) {
			return false
		}
	}

	return true
}

func DecideKeywords(title string) []string {
	var keywords []string
	for _, key := range keys {
		Next:
			for _, w := range key.Words {
				var find bool
				if w.Logic == "or" {
					find = orLogic(title, w)
				} else if w.Logic == "and" {
					find = andLogin(title, w)
				}
				if find {
					keywords = append(keywords, key.Keyword)
					continue Next
				}
			}
	}

	if len(keywords) == 0 {
		keywords = append(keywords, "其他")
	}
	return keywords
}

