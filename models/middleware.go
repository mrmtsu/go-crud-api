package models

import (
	"fmt"
	"go-rest-api/config"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Post struct {
	Model
	Title string `gorm:"size:255" json:"title"`
	Text  string `gorm:"size:255" json:"text"`
}

var Db *gorm.DB

func GetAllPosts(posts *[]Post) {
	Db.Find(&posts)
}

func GetSinglePost(post *Post, key string) {
	Db.First(&post, key)
}

func InsertPost(post *Post) {
	Db.NewRecord(post)
	Db.Create(&post)
}

func DeletePost(key string) {
	Db.Where("id = ?", key).Delete(&Post{})
}

func UpdatePost(post *Post, key string) {
	Db.Model(&post).Where("id = ?", key).Updates(
		map[string]interface{}{
			"title": post.Title,
			"text":  post.Text,
		},
	)
}

func init() {
	var err error
	dbConnectInfo := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName,
	)

	Db, err = gorm.Open(config.Config.DbDriverName, dbConnectInfo)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Successfully connect database..")
	}

	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Post{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfulle created table..")
	}
}
