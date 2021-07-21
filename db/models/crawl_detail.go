package models

import "snake/db"

type CrawlDetail struct {
	ID int `gorm:"primary_key"`
	Url string `gorm:"type:varchar(128);"`
	Title string `gorm:"type:varchar(256);"`
	Keyword string `gorm:"type:varchar(64);"`
	Timestamp int `gorm:"type:int(11);"`
}


func FindByDetailId(id int) CrawlDetail{
	var detail CrawlDetail
	db.DB.Where("id = ?", id).Find(&detail)
	return detail
}

func AddCrawlDetail(url string, title string, timestamp int) {
	var detail CrawlDetail
	db.DB.Where("url = ?", url).First(&detail)
	if detail.ID == 0 {
		db.DB.Create(&CrawlDetail{
			Url: url,
			Title: title,
			Timestamp: timestamp,
		})
	}
}

func AllCrawlDetails() []CrawlDetail {
	var details []CrawlDetail
	db.DB.Order("timestamp desc").Find(&details)
	return details
}