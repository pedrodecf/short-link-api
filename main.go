package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", "postgres://docker:docker@localhost:5432/shortlinks?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

var rdb *redis.Client

func initRedis() {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "docker",
		DB:       1,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS short_links (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			short_code TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

type link struct {
	Url       string `json:"url"`
	ShortCode string `json:"short_code"`
}

func postNewLink(c *gin.Context) {
	var newLink link

	if err := c.BindJSON(&newLink); err != nil {
		return
	}

	_, err := db.Exec("INSERT INTO short_links (url, short_code) VALUES ($1, $2)", newLink.Url, newLink.ShortCode)

	if err != nil {
		log.Fatal(err)
	}
}

func getLink(c *gin.Context) {
	shortCode := c.Param("short_code")

	var url string
	db.QueryRow("SELECT url FROM short_links WHERE short_code = $1", shortCode).Scan(&url)

	ctx := context.Background()
	rdb.ZIncrBy(ctx, "clicks", 1, shortCode)

	c.Redirect(301, url)
}

func getAllLinks() []link {
	rows, err := db.Query("SELECT url, short_code FROM short_links")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var links []link
	for rows.Next() {
		var l link
		err := rows.Scan(&l.Url, &l.ShortCode)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, l)
	}

	return links
}

func metrics(c *gin.Context) {
	result := rdb.ZRangeWithScores(context.Background(), "clicks", 0, -1)

	c.IndentedJSON(200, result.Val())
}

func main() {
	router := gin.Default()
	router.POST("/api/link", postNewLink)
	router.GET("/:short_code", getLink)
	router.GET("/api/metrics", metrics)
	router.GET("/api/links", func(c *gin.Context) {
		c.IndentedJSON(200, getAllLinks())
	})
	initDB()
	initRedis()
	createTable()
	router.Run("localhost:8080")
}
