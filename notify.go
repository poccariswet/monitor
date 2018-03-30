package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

	wg.Add(1)
	go ProcessHandle("./test")

	wg.Wait()
}

func ProcessHandle(bin string) {
	defer wg.Done()
	cmd := exec.Command(bin)
	cmd.Start()

	cmd.Wait()
	fmt.Println("cmd: ", cmd.Process.Pid)
	if cmd.ProcessState.Exited() != true {
		log.Print(bin, " is not exited")
	}

	if err := Notify(bin); err != nil {
		log.Fatal(err)
	}
}

func Notify(text string) error {
	msg := fmt.Sprintf("%s proceess が止まりました。", text)

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
