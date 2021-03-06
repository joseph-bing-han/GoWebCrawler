package model

import (
	"GoWebCrawler/src/utils/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var DB *gorm.DB

// Items [...]
type Item struct {
	gorm.Model
	Category   string `gorm:"column:category;type:varchar(100)" json:"category"`
	Website    string `gorm:"unique_index:item_production_id;column:website;type:varchar(20)" json:"website"`
	ProductID  string `gorm:"unique_index:item_production_id;column:product_id;type:varchar(255)" json:"product_id"`
	InternalID string `gorm:"column:internal_id;type:varchar(255)" json:"internal_id"`
	Title      string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	TitleZh    string `gorm:"column:title_zh;type:varchar(255);not null" json:"title_zh"`
	Unit       string `gorm:"column:unit;type:varchar(30)" json:"unit"`
	Url        string `gorm:"unique;column:url;type:varchar(255)" json:"url"`
	Image      string `gorm:"column:image;type:varchar(255)" json:"image"`
	Prices     []Price
}

// Prices [...]
type Price struct {
	ID        int64     `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"-"`
	ItemID    int64     `gorm:"index;column:item_id;type:bigint(20) unsigned" json:"item_id"`
	Items     Item      `gorm:"association_foreignkey:item_id;foreignkey:id" json:"items_list"`
	Branch    string    `gorm:"column:branch;type:varchar(255)" json:"branch"`
	Price     float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
}

// Sources [...]
type Source struct {
	gorm.Model
	Category string `gorm:"column:category;type:varchar(100)" json:"category"`
	Url      string `gorm:"column:url;type:varchar(255)" json:"url"`
	Active   bool   `gorm:"column:active;type:tinyint(1);not null" json:"active"`
}

func init() {

	dbHost := conf.Get("DB_SERVER", "127.0.0.1")
	dbPort := conf.Get("DB_PORT", "3306")
	dbName := conf.Get("DB_DATABASE", "crawler")
	dbUser := conf.Get("DB_USER", "root")
	dbPassword := conf.Get("DB_PASSWORD", "root")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysql, err := gorm.Open("mysql", dsn)
	if err == nil {
		DB = mysql
	}
}
