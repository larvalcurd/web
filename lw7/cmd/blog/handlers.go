package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	FeaturedPosts []featuredPostData
	RecentPosts   []recentPostData
}

type featuredPostData struct {
	Title     string `db:"title"`
	SubTitle  string `db:"subtitle"`
	Author    string `db:"author"`
	AuthorImg string `db:"author_url"`
	PubDate   string `db:"publish_date"`
	PostImg   string `db:"image_url"`
}

type recentPostData struct {
	Title         string `db:"title"`
	SubTitle      string `db:"subtitle"`
	RecentPostImg string `db:"image_url"`
	Author        string `db:"author"`
	AuthorImg     string `db:"author_url"`
	PubDate       string `db:"publish_date"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/index.html")

		if err != nil {
			http.Error(w, "Internal Server Error 1", 500)
			log.Println(err.Error())
			return
		}

		featured, err := featuredPosts(db)

		if err != nil {
			http.Error(w, "Internal Server Error featured", 500)
			log.Fatal(err)
			return
		}

		recent, err := recentPosts(db)

		if err != nil {
			http.Error(w, "Internal Server Error recent", 500)
			log.Fatal(err)
			return
		}

		data := indexPage{
			FeaturedPosts: featured,
			RecentPosts:   recent,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error 2", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")

	}
}

func featuredPosts(client *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			image_url,
			image_url,
			author,
			author_url,
			publish_date
		FROM 
			post
		WHERE featured = 1
	`

	var posts []featuredPostData

	err := client.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                     // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}

func recentPosts(client *sqlx.DB) ([]recentPostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			image_url,
			author,
			author_url,
			publish_date
		FROM 
			post
		WHERE featured = 0
	`

	var posts []recentPostData

	err := client.Select(&posts, query)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
