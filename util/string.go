package util

import (
	"regexp"
	"strings"
)

func DeleteExtractSpace(s string) string {
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

// 返回string 长度小于l的rune
func Substring(source string, l int) string {
	var r = []rune(source)
	length := len(r)

	var substring = ""
	for i := 0; i < length && len(substring) < l; i++ {
		if len(substring) > (l - len(string(r[i]))) {
			break
		}
		substring += string(r[i])
	}

	return substring
}

func TrimEmoji(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
		} else {
			ret += string(rs[i])
		}
	}
	return ret
}