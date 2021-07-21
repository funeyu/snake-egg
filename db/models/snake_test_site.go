package models

import "snake/db"

type SnakeTestSite struct {
	ID int `gorm:"primary_key"`
	Url string `gorm:"type:varchar(128);"`
	Host string `gorm:"type:varchar(64);"`
	Speed int `gorm:"type:smallint;"`
	Title string `gorm:"type:varchar(128);"`
	Keywords string `gorm:"type:varchar(256);"`
	Favicon string `gorm:"type:varchar(128);"`
	Score int `gorm:"type:smallint;"`
}


func AllSites() []SnakeTestSite {
	var sites []SnakeTestSite
	db.DB.Find(&sites)
	return sites
}