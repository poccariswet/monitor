package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

const (
	endpoint = "https://notify-api.line.me/api/notify"
	msg_good = "running write process"
	msg_bad  = "stopped write process\nPlease check write server"
)

var (
	accessToken = os.Getenv("NOTIFY_TOKEN")
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please set 'monitor <pid>'")
		os.Exit(1)
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			pidInfo, _ := ps.FindProcess(pid)
			if pidInfo == nil {
				Notify(msg_bad)
				os.Exit(1)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		pidInfo, _ := ps.FindProcess(pid)
		if pidInfo == nil {
			Notify(msg_bad)
			os.Exit(1)
		}
		Notify(msg_good)
		time.Sleep(600 * time.Second)
		//		time.Sleep(12 * time.Hour)
	}

}

// send message using notify
func Notify(msg string) error {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return err
	}

	c := &http.Client{}
	form := url.Values{}
	form.Add("message", msg)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
