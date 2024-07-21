package main

import (
	"context"
	"fmt"
	"go-crawler/internal/constant"
	"go-crawler/internal/models"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(os.Getenv("DATABASE_NAME"))
	c := colly.NewCollector()

	var news = []models.News{}
	var documents []interface{}

	c.OnHTML("article.articles--iridescent-list--item", func(e *colly.HTMLElement) {
		title := e.ChildText("header h4.articles--iridescent-list--text-item__title a span.articles--iridescent-list--text-item__title-link-text")
		url := e.ChildAttr("header h4.articles--iridescent-list--text-item__title a.ui--a", "href")

		if title != "" || url != "" {
			fmt.Printf("Title: %s\n", title)
			fmt.Printf("URL: %s\n", url)
			news = append(news, models.News{
				Title: title,
				Url:   url,
			})
		}

	})
	c.Visit(constant.LIPUTAN6_URL)

	for _, val := range news {
		documents = append(documents, val)
	}

	if len(news) > 0 {
		res, err := db.Collection(constant.LIPUTAN6_COLLECTION).InsertMany(context.TODO(), documents)
		if err != nil {
			log.Println(err.Error())
			log.Fatal(err)
		}
		fmt.Printf("%d documents inserted with IDs:\n", len(res.InsertedIDs))
		for _, id := range res.InsertedIDs {
			fmt.Printf("\t%s\n", id)
		}
	}

}
