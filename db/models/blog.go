package models

import (
	"snake/db"
	"snake/util"
)

type Blog struct {
	ID int `gorm:"primary_key"`
	GitLogin string `gorm:"type:varchar(32);"`
	GitFollowers int `gorm:"type:int;"`
	Domain string `gorm:"type:varchar(32);"`
	Schema int `gorm:"type:varchar(32);"`
	Path string `gorm:"type:varchar(32);"`
	Speed int `gorm:"type:int;"`
	ArticleNum int `gorm:"type:int;"`
	Favicon string `gorm:"type:varchar(128);"`
	Title string `gorm:"type:varchar(32);"`
	Keywords string `gorm:"type:varchar(128);"`
	Author string `gorm:"type:varchar(32);"`
	Lang uint8  `gorm:"type:tinyint(1);"`
	LastBlogUrl string `gorm:"type:varchar(128);"`
	RankId int `gorm:"type:int"`
	SubRankId int `gorm:"type:int"`
	FileId uint8 `gorm:"type:tinyint;"`
	FileFilled uint8 `gorm:"type:tinyint;"`
	BlockNum int `gorm:"type:int;"`
}

func (b *Blog) FormUrl() string {
	schema := "http://"
	if b.Schema == 1 {
		schema = "https://"
	}

	return schema + b.Domain + b.Path
}

func (b *Blog) UpdateArticleNum(n int) {
	b.ArticleNum = n
	db.DB.Save(b)
}

func (b *Blog) UpdateFile(fileId uint8, blockNum int) {
	db.DB.Model(&b).Update(map[string]interface{}{"file_id": fileId, "block_num": blockNum, "file_filled": 1})
}

func AppendOrNotBlog(b *Blog) {
	var find Blog
	db.DB.First(&find, "domain=?", b.Domain)

	if find.ID == 0 {
		b.Keywords = util.TrimEmoji(b.Keywords)
		db.DB.Create(b)
	} else {
		db.DB.Model(&find).Updates(map[string]interface{}{
			"favicon": b.Favicon,
			"git_followers": b.GitFollowers,
		})
	}
}

func UpdateOrAddBlog(b Blog) {
	b.Keywords = util.TrimEmoji(b.Keywords)
	if b.ID != 0 {
		db.DB.Model(&b).Updates(b)
	} else {
		var bb Blog
		db.DB.Where("domain=?", b.Domain).First(&bb)
		if bb.ID != 0 {
			db.DB.Model(&bb).Update(b)
		} else {
			db.DB.Create(&b)
		}
	}
}

func AllBlogs() []Blog {
	var blogs []Blog
	db.DB.Order("rank_id asc, sub_rank_id").Find(&blogs)
	//db.DB.Where("git_login=?", "TechJene").Find(&blogs)
	return blogs
}
func AllBlogsWithArticle() []Blog {
	var blogs []Blog
	db.DB.Where("article_num != ?", 0).Find(&blogs)
	return blogs
}

func BlogsRankId(rankid int) Blog {
	var blogs []Blog
	db.DB.Where("rank_id=?", rankid).Order("sub_rank_id desc").Limit(1).First(&blogs)
	if len(blogs) == 1 {
		return blogs[0]
	}
	return Blog{}
}

func CreateBlog(blog Blog) {
	blog.Keywords = util.TrimEmoji(blog.Keywords)
	db.DB.Create(&blog)
}

func FindBlog(login string) bool {
	var blog Blog
	db.DB.Where("git_login=?", login).First(&blog)
	if blog.ID != 0 {
		return true
	}
	return false
}
