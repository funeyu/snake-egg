package jieba

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/**
整个结巴分词用到的预料库，可能不适合我们博客，后续优化 该语料库
*/

var di *dict

func init() {
	di = load("./jieba/resource/dict.txt.big")
}

type path struct {
	idx  int
	freq float64
}

func find_max(paths []path) *path {
	var p *path
	if len(paths) == 0 {
		return p
	}
	p = &paths[0]
	for _, pp := range paths {
		if pp.freq > p.freq {
			p = &pp
		}
	}

	return p
}


// 这里注意是要将中英文的混合给正确分离, 如 "Linux I/O 多路复用"
func split_words(sentence string) []string {
	splits := make([][]string, 1)

	for _, s := range sentence {
		if !unicode.Is(unicode.Han, s) {
			if unicode.IsSpace(s) {
				s := make([]string, 1)
				splits = append(splits, s)
			} else {
				l := splits[len(splits) - 1]
				if len(l) == 1 {
					r, _ := utf8.DecodeRuneInString(l[0])
					if unicode.Is(unicode.Han, r) { // 英文处于中文之间 如 "分布式ID生成算法"
						splits = append(splits, []string {string(s)})
						continue
					}
				}

				splits[len(splits) - 1] = append(splits[len(splits) - 1], string(s))
			}
		} else {
			splits = append(splits, []string{string(s)})
		}
	}

	var words []string
	for _, s := range splits {
		if len(s) == 1 {
			words = append(words, strings.ToLower(s[0]))
		} else if len(s) > 1 {
			words = append(words, strings.ToLower(strings.Join(s, "")))
		}
	}
	return words
}

/**
根据用户要切分的语句形成有向图
比如"化身成龙" 的dag为：
{
	0: [0, 1],
	1: [1],
	2: [2, 3],
	3: [3]
}
*/
func generate_dag(sentence string) map[int][]int {
	words := split_words(sentence)
	lens := len(words)
	DAG := make(map[int][]int, lens)

	for i, _ := range words {
		DAG[i] = make([]int, 0)
		for j := i; j < lens; j++ {
			ss := strings.Join(words[i:j+1], "")
			if di.contains(ss) {
				DAG[i] = append(DAG[i], j)
			}
		}
	}

	return DAG
}

func Cut(sentence string) []string {
	dag := generate_dag(sentence)
	words := split_words(sentence)
	cap := len(dag)

	routes := make([]path, cap+1)
	routes[cap] = path{
		idx:  cap,
		freq: 0,
	}

	for i := cap - 1; i >= 0; i-- { //计算每个routes值
		paths := make([]path, len(dag[i]))
		for ii, d := range dag[i] {
			word := strings.Join(words[i:d+1], "")
			paths[ii] = path{
				idx:  d,
				freq: di.freq(word) + routes[d+1].freq,
			}
		}

		max_path := find_max(paths)
		if max_path !=nil {
			routes[i] = path{
				idx:  max_path.idx,
				freq: max_path.freq,
			}
		} else {
			routes[i] = path {
				idx: i,
				freq: 0,
			}
		}

	}

	var cuts []string
	for c := 0; c < cap; {
		r := routes[c]
		word := strings.Join(words[c:r.idx+1], "")
		cuts = append(cuts, trim(word)...)
		c = r.idx + 1
	}

	return cuts
}
