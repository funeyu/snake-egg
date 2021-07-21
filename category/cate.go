package category

var keywords = map[string][]string {
	"程序技术博客": []string{"python", "前端", "numpy", "全栈", "flutter"},
	"在线工具": []string{"pdf", "生成工具", "生成工具", "工具", "各类工具", "工具箱", "转换器"},
	"视频工具": []string{"视频工具"},
	"美女直播": []string{"美女视频", "花椒直播"},
	"有趣手游网": []string{"手游"},
	"性感美女": []string{"美女", "美女动态图", "美女写真"},
	"热门新闻": []string{"全网资讯"},
	"装逼": []string{"装逼", "在线装逼"},
	"网址导航": []string{"网址导航", "网站导航"},
	"文艺心": []string{"文艺", "金句"},
	"视频短片": []string{"微电影", "视频", "微电影", "视频教育"},
	"资源神器": []string{"书格", "电子书", "资源神器", "资源搜索", "免费下载", "百度网盘", "网盘搜索", "资料网"},
}


func _contains(words []string, target string) bool {
	for _, w := range words {
		if w == target {
			return true
		}
	}
	return false
}

func Sort(words []string) string {
	for k, v := range keywords {
		for _, w :=  range v {
			if _contains(words, w) {
				return k
			}
		}
	}

	return ""
}
