package utils

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func dotENV(key string) string {
	err := godotenv.Load(".env")
	exit(err)
	return os.Getenv(key)
}

func dotENVInt(key string) int {
	err := godotenv.Load(".env")
	exit(err)
	i, err := strconv.Atoi(os.Getenv(key))
	exit(err)
	return i
}

type Server struct{}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Data []Item `xml:"item" json:"item"`
}

type Item struct {
	Title       string  `xml:"title" json:"Title"`
	Link        string  `xml:"guid" json:"Link"`
	Description string  `xml:"description" json:"Description"`
	Image       []image `xml:"content" json:"Image"`
}

type image struct {
	URL string `xml:"url,attr" json:"URL"`
}

type Store map[string][]Item

var Client = redis.NewClient(&redis.Options{
	Addr:     dotENV("REDIS_ADDR"),
	Password: dotENV("REDIS_PASSWORD"),
	DB:       dotENVInt("REDIS_DB"),
})

var Ticker *time.Ticker

var urls []string = []string{
	"https://moxie.foxnews.com/google-publisher/world.xml",
	"https://www.theguardian.com/uk/rss",
	"https://rss.nytimes.com/services/xml/rss/nyt/world.xml",
	"https://feeds.feedburner.com/ndtvnews-world-news",
	"https://www.cgtn.com/subscribe/rss/section/world.xml",
	"https://www.hindustantimes.com/feeds/rss/videos/world-news/rssfeed.xml",
}

var mutex sync.Mutex = sync.Mutex{}
