package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ps "github.com/mitchellh/go-ps"
)

const (
	endpoint = "https://hooks.slack.com/services/"
	msg_good = "running write process"
	msg_bad  = "stopped write process\nPlease check write server"
)

var (
	accessToken = os.Getenv("INCOMMING_TOKEN")
	quit        = make(chan struct{})
	thour       = make(chan int)
)

// send message using notify
func Notify(msg string) error {
	url := endpoint + accessToken
	data := `{"text":"` + msg + `"}`
	body := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
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
