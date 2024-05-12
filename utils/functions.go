package utils

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func ExtractTextBetweenPTags(input string) string {
	input = "<p>" + input + "</p>"
	re := regexp.MustCompile(`<p>([^<]+)</p>`)
	matches := re.FindAllStringSubmatch(input, -1)
	var result []string
	for _, match := range matches {
		result = append(result, match[1])
	}
	return strings.Join(result, " ")
}

func ImagePrint(Image []image) string {
	if len(Image) == 0 {
		return "No Image."
	}
	if len(Image) == 1 {
		return Image[0].URL
	}
	return Image[1].URL
}

func exit(err error) {
	if err != nil && err != redis.Nil {
		fmt.Printf("Error : %v\n", err)
		os.Exit(1)
	}
}

func XMLParser(url string) []Item {
	res, err := http.Get(url)
	exit(err)
	defer res.Body.Close()

	xmlData, err := io.ReadAll(res.Body)
	exit(err)

	var tmp RSS
	err = xml.Unmarshal(xmlData, &tmp)
	exit(err)

	feed := make(chan Item)
	for _, dpoint := range tmp.Channel.Data {
		go func(item Item) {
			item.Description = ExtractTextBetweenPTags(item.Description)
			item.Image = []image{
				{URL: ImagePrint(item.Image)},
			}
			feed <- item
		}(dpoint)
	}

	var dataAfterParsed []Item = make([]Item, len(tmp.Channel.Data))
	for i := range tmp.Channel.Data {
		dataAfterParsed[i] = <-feed
	}

	return dataAfterParsed
}

func JSONify(data []Item) string {
	jData, err := json.Marshal(data)
	exit(err)
	return string(jData)
}

func Randomizer(data *[]Item) {
	for i := len(*data) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		(*data)[i], (*data)[j] = (*data)[j], (*data)[i]
	}
}

func GetData() []Item {
	val, err := Client.Get(context.Background(), "data").Result()
	exit(err)

	if err != redis.Nil {
		// Value is present
		return deserialize(val)
	}

	var dataAfterParsed []Item
	var feed = make(chan []Item)
	for _, url := range urls {
		go func(url string) {
			feed <- XMLParser(url)
		}(url)
	}
	dataAfterParsed = fillHelper(feed)
	Randomizer(&dataAfterParsed)

	err = Client.Set(context.Background(), "data", serialize(dataAfterParsed), 1*time.Hour).Err()
	exit(err)
	return dataAfterParsed
}

func serialize(dataAfterParsed []Item) string {
	jsonData, err := json.Marshal(dataAfterParsed)
	exit(err)
	return string(jsonData)
}

func fillHelper(feed chan []Item) []Item {
	mutex.Lock()
	defer mutex.Unlock()
	var dataAfterParsed []Item
	for range urls {
		dataAfterParsed = append(dataAfterParsed, <-feed...)
	}
	return dataAfterParsed
}

func Revalidate() {
	Ticker = time.NewTicker(50 * time.Minute)
	go func() {
		for {
			fmt.Println("Revalidating cache....")
			GetData()
			<-Ticker.C
		}
	}()
}

func deserialize(val string) (dataTobeSent []Item) {
	e := json.Unmarshal([]byte(val), &dataTobeSent)
	exit(e)
	return
}
