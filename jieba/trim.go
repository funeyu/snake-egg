package jieba

import (
	"github.com/araddon/dateparse"
	"regexp"
	"strings"
)

// 是否为无意义的string，如纯数字，纯时间string
func noMeaning(s string) bool {
	_, err := dateparse.ParseAny(s)
	if err == nil {
		return true
	}

	number := regexp.MustCompile(`^\d+$`)
	if number.MatchString(s) {
		return true
	}
	return false
}

func trim(s string) []string {
	var res []string
	reg := regexp.MustCompile(`(\(|\)|\[|\]|【|】|\?|？|\:|「|」|》|《|\>|\<|"|“)`)
	s = reg.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	if strings.Contains(s, "http") {
		return nil
	}
	if noMeaning(s) {
		return nil
	}
	if strings.Contains(s, "+") {
		res = append(res, strings.Split(s, "+")...)
	}
	if len(res) == 0 && strings.Contains(s, "-") {
		res = append(res, strings.Split(s, "-")...)
	}
	if len(res) == 0 {
		res = append(res, s)
	}
	return res
}
