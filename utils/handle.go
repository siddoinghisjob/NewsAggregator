package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	path := r.URL.Path
	pattern := `^\/\d+$`
	matched, err := regexp.MatchString(pattern, path)
	exit(err)

	if !matched {
		err := fmt.Errorf("url format incorrect")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	page, err := strconv.Atoi(path[1:])
	exit(err)

	xmlData := GetData()
	var output []Item
	start := (page - 1) * 10

	if len(xmlData) > (page-1)*10 {
		if len(xmlData) < page*10 {
			output = xmlData[start:]
		} else {
			end := page * 10
			output = xmlData[start:end]
		}
	} else {
		err := fmt.Errorf("url format incorrect")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	jsonData := JSONify(output)

	fmt.Fprint(w, jsonData)
}
