package main

import (
	"fmt"
	"github.com/funeyu/smallfiles"
	"snake/db/models"
	"snake/indexer"
	"snake/jieba"

	"snake/util"

	"strconv"
	//"strings"
	"testing"
)

func TestFilterData(t *testing.T) {
	engine := indexer.Init()
	sf := smallfiles.Open("/Users/fuheyu/snake/snake-crawl/data/", ArticleFormat)
	b, _ := sf.GetBlock(0, 1)
	for _, d := range b.Datas {
		a := d.(*Article)
		a.Text = util.Substring(a.Text, 250)
		keys := jieba.Cut(a.Text)
		id := strconv.FormatUint(uint64(a.Id), 10)
		d := indexer.Doc{
			Id:                 id,
			DocId:              a.Id,
			Url:                a.Href,
			Lang:				a.Lang,
			Title:              a.Text,
			TimeStamp:          uint32(a.TimeStamp),
			Favicon:            a.Favicon,
			TitleOrDescription: a.KeywordOrDescription,
			Star:               0,
			IsTop5:             false,
		}
		fmt.Println("artile", d)
		engine.AddDocOrderly(&d, keys)
	}
	ids := engine.Search([]string{"typesript"}, false)
	fmt.Println("ids", ids)
}

func TestFilterData2(t *testing.T) {
	//engine := indexer.Init()
	//store := store.InitBadger("./badger")
	blogs := models.AllBlogs()
	total := 0
	sf := smallfiles.Open("/Users/fuheyu/snake/snake-crawl/data/", ArticleFormat)
	for _, blog := range blogs {
		if blog.FileFilled ==0 {
			continue
		}
		b, e  := sf.GetBlock(blog.FileId, blog.BlockNum)
		if e == nil {
			if b == nil {
				fmt.Println("b.nil", blog)
				continue
			}
			for i, d := range b.Datas {
				if d == nil {
					continue
				}
				a := d.(*Article)
				a.Text = util.Substring(a.Text, 250)
				a.KeywordOrDescription = util.Substring(a.KeywordOrDescription, 250)
				docId, _ := indexer.GeneDocId(uint32(blog.RankId), uint32(blog.SubRankId), uint16(i))
				a.Id = docId

				////keys := jieba.Cut(a.Text)
				//id := strconv.FormatUint(uint64(a.Id), 10)
				//d := indexer.Doc{
				//	Id:                 id,
				//	DocId:              a.Id,
				//	Url:                a.Href,
				//	Lang:				a.Lang,
				//	Title:              a.Text,
				//	TimeStamp:          uint32(a.TimeStamp),
				//	Favicon:            a.Favicon,
				//	TitleOrDescription: a.KeywordOrDescription,
				//	Star:               0,
				//	IsTop5:             false,
				//}
				//if strings.Contains(a.Text, "golang") {
				//	total = total + 1
				//}
				////engine.AddDocOrderly(&d, keys)
				//store.Add(d)
			}
			sf.RefillDatas(b.Datas, blog.FileId, blog.BlockNum)
		}
	}
	fmt.Println("total", total)
	//engine.FlushDisk()
}

func TestFilterData3(t *testing.T) {
	sf := smallfiles.Open("/Users/fuheyu/snake/snake-crawl/data/", ArticleFormat)
	b, _ := sf.GetBlock(51, 86)
	for _, d := range b.Datas {
		a := d.(*Article)

		fmt.Println(a.Id)
	}
}