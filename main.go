package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/siddoinghisjob/news-aggregator/utils"
)

func main() {
	defer utils.Client.Close()
	utils.Revalidate()
	fmt.Println("Listening to port 8080....")
	log.Fatal(http.ListenAndServe(":8080", &utils.Server{}))
}
