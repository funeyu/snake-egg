package db

type RealBlog struct {
	ID int `gorm:"primary_key"`
	Domain string `gorm:"type:varchar(256);"`
	Speed int `gorm:"type:int;"`
	Star int `gorm:"type:int"`
}

func All() []RealBlog {
	var blogs []RealBlog
	DB.Find(&blogs)
	return blogs
}


func DeleteOnFindDomain (domain string) {
	var blog RealBlog
	DB.First(&blog, "domain=?", domain)

	if blog.ID != 0 {
		DB.Delete(&blog)
	}
}

func AppendOrNot(domain string) {
	var blog RealBlog
	DB.First(&blog, "domain=?", domain)

	if blog.ID == 0 {
		DB.Create(&RealBlog{
			Domain: domain,
		})
	}
}