package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Link struct {
	Slug string
	Data string
	gorm.Model
}

func initDb() *gorm.DB {
	fmt.Println("Initializing the DB")

	dsn := "host=postgres user=crs password=crs dbname=crs-pics port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to db")
		panic("DB connection failed")
	}

	fmt.Println("Connection to Postgres established!")

	err = database.AutoMigrate(&Link{})
	if err != nil {
		fmt.Println("Migration failed")
		panic("Migration failed")
	}

	return database
}

type db struct {
	*gorm.DB
}

type DB interface {
	GetAllLinks() []Link
	GetLinkBySlug(slug string) string
	AddLink(link Link) *Link
}

func NewDB() DB {
	database := initDb()
	return &db{
		database,
	}
}

func (db *db) GetAllLinks() []Link {
	var links []Link
	db.Find(&links)

	return links
}

func (db *db) GetLinkBySlug(slug string) string {
	var link Link
	db.Where("slug = ?", slug).Take(&link)

	return link.Data
}

func (db *db) AddLink(link Link) *Link {
	result := db.Clauses(clause.Returning{}).Create(&link)
	fmt.Println(result.RowsAffected)

	return &link
}
