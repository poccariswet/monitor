package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	wg          sync.WaitGroup
	accessToken = os.Getenv("NOTIFY_TOKEN")
)

const (
	endpoint = "https://notify-api.line.me/api/notify"
)

func main() {
	msg := "write process"

	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	c := &http.Client{}
	form := url.Values{}
	form.Add("message", msg)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
