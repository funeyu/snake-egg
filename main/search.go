package main

import (
	"bufio"
	"fmt"
	"github.com/funeyu/smallfiles"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"snake/db/models"
	"snake/implement"
	"snake/indexer"
	"snake/jieba"

	//"snake/jieba"
	search "snake/proto"
	"snake/store"
	"snake/util"
	"strconv"
	"strings"
)

func scanFile(fileDir string, m map[string]string) {
	//fileDir := "./main/data3"
	files, err:= ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("err", err)
	}
	re := regexp.MustCompile(`:[\s]+`)

	for _, file := range files {
		f, _ := os.Open(fileDir + file.Name())
		scan := bufio.NewScanner(f)
		scan.Split(bufio.ScanLines)

		for scan.Scan() {
			t := scan.Text()
			ss := re.Split(t, -1)
			if len(ss) != 2 {
				continue
			}
			m[ss[0]] = ss[1]
		}
	}
}

func FilterData() map[string]string {
	res := make(map[string] string, 0)
	scanFile("./data/", res)
	scanFile("./data2/",res)
	return res
}

type Link struct {
	Text string `json: "text"`
	Href string `json: "href"`
	TimeStamp int `json:"time_stamp"`
}

type Article struct {
	Link
	Lang uint8 `json:"lang"`
	Id indexer.DocId `json:"id"`
	Favicon string `json:"favicon"`
	Title string `json:"title"`
	KeywordOrDescription string `json:"keyword_or_description"`
}

//type SmallData interface { // 标识存储的最下数据单元，如存取的一条文章信息
//	Size() uint32
//	Serialize() []byte
//}

func (a *Article) Serialize() []byte {
	id := strconv.FormatUint(uint64(a.Id), 10)
	timestamp := strconv.FormatInt(int64(a.TimeStamp), 10)
	lang := strconv.FormatUint(uint64(a.Lang), 10)
	// 每个文章存储在smallfile中的字段如下， 以'##'分隔
	str := strings.Join([]string{id, a.Text, a.Href ,timestamp,  a.Favicon, a.KeywordOrDescription, lang}, "##")
	return []byte(str)
}

func (a *Article) Size() uint32 {
	id := strconv.FormatUint(uint64(a.Id), 10)
	timestamp := strconv.FormatInt(int64(a.TimeStamp), 10)
	lang := strconv.FormatUint(uint64(a.Lang), 10)
	// 每个字段string的长度+ '##'的长度
	total := len(id) + len(a.Text) + len(a.Href) +  len(timestamp) + len(a.KeywordOrDescription) + len(lang)  + 6 * 2
	return uint32(total)
}

func ArticleFormat(bytes []byte) smallfiles.SmallData {
	str := string(bytes)
	ss := strings.Split(str, "##")
	id,_ := strconv.ParseUint(ss[0], 10, 64)
	timeStamp, _ := strconv.ParseUint(ss[3], 10, 64)
	lang, _ := strconv.ParseUint(ss[5], 10, 64)

	return &Article{
		Link:                 Link{ Text: ss[1], Href:ss[2], TimeStamp: int(timeStamp)},
		Lang:                 uint8(lang),
		Id:                   indexer.DocId(id),
		Title:                ss[1],
		KeywordOrDescription: ss[4],
	}
}

func loadData(store store.Store, engine *indexer.Indexer) {
	blogs := models.AllBlogs()
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
			for _, d := range b.Datas {
				if d == nil {
					continue
				}
				a := d.(*Article)
				a.Text = util.Substring(a.Text, 250)
				a.KeywordOrDescription = util.Substring(a.KeywordOrDescription, 250)

				if strings.Contains(a.Text, "-color-scheme") {
					fmt.Println("a.Text", a.Text, a)
				}
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
				engine.AddDocOrderly(&d, keys)
			}
		}
	}
	engine.FlushDisk()
}


func main() {
	store := store.InitBadger("./badger")
	engine := indexer.Init()
	// listen on message queue
	//consumer, _ := nsq.CreateDocConsumer("egg")
	//consumer.Listen(func(d indexer.Doc) error {
	//	fmt.Println("dynamic doc", d)
	//	title := util.DeleteExtractSpace(d.Title)
	//	keywords := jieba.Cut(d.Title)
	//	engine.AddDocDynamic(&d, keywords)
	//	doc := indexer.Doc{
	//		Id:                 d.Id,
	//		DocId:              d.DocId,
	//		Url:                d.Url,
	//		Lang: 				d.Lang,
	//		Title:             	title,
	//		TimeStamp:          d.TimeStamp,
	//		Favicon:            d.Favicon,
	//		TitleOrDescription: d.TitleOrDescription,
	//		Star:               0,
	//		IsTop5:             false,
	//	}
	//	store.Add(doc)
	//	return nil
	//})
	//consumer.Init()
	loadData(store, engine)
	//engine.LoadFromDisk()
	ss := &implement.SearchServer{
		Store:   store,
		Indexer: engine,
	}

	starServer := implement.InitStarServer(&engine.StarRating)

	addr := ":50051"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	search.RegisterSearcherServer(s, ss)
	search.RegisterStarerServer(s, starServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
