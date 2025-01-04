package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Link struct {
	Slug      string `gorm:"default:random_string(7)"`
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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

	// db function to create a random string (for slug generation)
	dbFuncSql := `
		CREATE OR REPLACE FUNCTION random_string(length INT)
		RETURNS TEXT AS $$
		DECLARE
			characters TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
			characters_length INT := LENGTH(characters);
			random_string TEXT := '';
			i INT;
		BEGIN
			IF length <= 0 THEN
				RAISE EXCEPTION 'Length must be a positive integer';
			END IF;
		
			FOR i IN 1..length LOOP
				random_string := random_string || substr(characters, floor(random() * characters_length + 1)::int, 1);
			END LOOP;
		
			RETURN random_string;
		END;
		$$ LANGUAGE plpgsql;
	`

	database.Exec(dbFuncSql)

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
