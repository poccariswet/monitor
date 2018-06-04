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

	ps "github.com/mitchellh/go-ps"
)

const (
	endpoint = "https://notify-api.line.me/api/notify"
	msg_good = "running write process"
	msg_bad  = "stopped write process\nPlease check write server"
)

var (
	accessToken = os.Getenv("NOTIFY_TOKEN")
	quit        = make(chan struct{})
	thour       = make(chan int)
)

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
		t1 := time.NewTicker(5 * time.Second)
		t2 := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-t1.C:
				pidinfo, _ := ps.FindProcess(pid)
				if pidinfo == nil {
					Notify(msg_bad)
					WriteServerDown()
					close(quit)
				}
			case <-t2.C:
				thour <- time.Now().Hour()
			}
		}
		t1.Stop()
		t2.Stop()
	}()

	for {
		select {
		case <-quit: //error処理
			fmt.Fprint(os.Stderr, "server down")
			os.Exit(1)
		case h := <-thour:
			if h == 11 || h == 17 {
				Notify(msg_good)
			}
		}
	}
}
