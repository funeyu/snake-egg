package models

import "snake/db"

type CrawlSite struct {
	ID int `gorm:"primary_key"`
	Url string `gorm:"type:varchar(128);"`
	KeywordsTitle string `gorm:"type:varchar(256);"`
	NeedCrawl uint8 `gorm:"type:tinyint(1);"`
	CrawlKeyword string `gorm:"type:varchar(64);"`
	SmallfileInfo string `gorm:"type:varchar(128);"`
}

func AllCrawlSites() []CrawlSite{
	var sites []CrawlSite
	db.DB.Find(&sites)
	return sites
}
