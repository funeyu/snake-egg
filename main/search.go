package main

import (
	"fmt"
	"github.com/funeyu/snakedocid"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"snake/db/models"
	"snake/implement"
	"snake/indexer"
	"snake/smallfiles"

	search "snake/proto"

	"snake/store"
	"snake/util"
	"strconv"
	"strings"
)


type Link struct {
	Text string `json: "text"`
	Href string `json: "href"`
	TimeStamp int `json:"time_stamp"`
}

type Article struct {
	Link
	Id snakedocid.DocId `json:"id"`
	Title string `json:"title"`
	TimeStamp uint64 `json:"timestamp"`
}

type SmallData struct {
	Article
	IsFirstData bool
}

func (d *SmallData) Serialize() []byte {
	if d.IsFirstData { // 首个数据只存储domain信息
		domain := d.Article.Domain()
		str := strings.Join([]string{domain}, "##")
		return []byte(str)
	}
	a := d.Article
	id := strconv.FormatUint(uint64(a.Id), 10)
	timestamp := strconv.FormatInt(int64(a.TimeStamp), 10)
	// 每个文章存储在smallfile中的字段如下， 以'##'分隔
	str := strings.Join([]string{id, a.Text, a.Path() ,timestamp}, "##")
	return []byte(str)
}

func (a *Article) Path() string {
	href := a.Link.Href
	domain := util.Domain(href)
	return strings.Replace(href, domain, "", -1)
}

func (a *Article) Domain() string {
	href := a.Link.Href
	return util.Domain(href)
}

func (d *SmallData) Size() uint32 {
	a := d.Article
	if d.IsFirstData {
		total := len(a.Domain())
		return uint32(total)
	}

	id := strconv.FormatUint(uint64(a.Id), 10)
	timestamp := strconv.FormatInt(int64(a.TimeStamp), 10)
	// 每个字段string的长度+ '##'的长度
	total := len(id) + len(a.Text) + len(a.Path()) +  len(timestamp) + 3 * 2
	return uint32(total)
}

func ArticleFormat(bytes []byte) smallfiles.SmallData {
	str := string(bytes)
	ss := strings.Split(str, "##")
	if len(ss) < 3 {
		domain := ss[0]
		a := Article{
			Link: Link{Href: domain},
		}
		return &SmallData{
			Article: a,
			IsFirstData: true,
		}
	}
	id,_ := strconv.ParseUint(ss[0], 10, 64)
	timeStamp, _ := strconv.ParseUint(ss[3], 10, 64)

	a := Article{
		Link:                 Link{ Text: ss[1], Href:ss[2], TimeStamp: int(timeStamp)},
		Id:                   snakedocid.DocId(id),
		Title:                ss[1],
		TimeStamp: timeStamp,
	}
	return &SmallData{
		Article: a,
		IsFirstData: false,
	}
}


type Handler func(d indexer.Doc)

func iterateBlogs(handle Handler ) {
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
				a := d.(*SmallData)
				a.Text = util.Substring(a.Text, 250)

				if strings.Contains(a.Text, "-color-scheme") {
					fmt.Println("a.Text", a.Text, a)
				}
				id := strconv.FormatUint(uint64(a.Id), 10)
				d := indexer.Doc{
					Id:                 id,
					Url:                a.Href,
					Title:              a.Text,
					TimeStamp:          uint32(a.TimeStamp),
				}
				handle(d)
			}
		}
	}
}

func loadData(store store.Store, engine *indexer.Indexer) {
	//iterateBlogs(func(d indexer.Doc) { // 生成一遍索引
	//	keys := jieba.Cut(d.Title)
	//	engine.AddDocOrderly(&d, keys)
	//})
	//engine.FlushDisk()
	//
	//iterateBlogs(func(d indexer.Doc) { // 生成一遍badger文件
	//	store.Add(d)
	//})
	//
	//engine.FlushDisk()

	store.ForEach(func(id int, keyword string) error {
		engine.AddId(indexer.SnakeId{
			Id: id,
		}, keyword)
		return nil
	})
}

func Refresh() {
	//store := store.InitBadger("./badger")
	store := store.InitDBStore()
	engine := indexer.Init()
	loadData(store, engine)
	//engine.LoadFromDisk()
	ss := implement.SearchServer{
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
	search.RegisterStarerServer(s, starServer)
	search.RegisterSearcherServer(s, &ss)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	c := cron.New()
	c.AddFunc("0 0 0 * * *", func() {
		Refresh()
	})
}
