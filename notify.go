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
	"time"
)

var (
	wg          sync.WaitGroup
	accessToken = os.Getenv("NOTIFY_TOKEN")
	con_text    = fmt.Sprintf("定時報告: 実行中(%s)", os.Getenv("HOME"))
)

const (
	endpoint = "https://notify-api.line.me/api/notify"
	duration = 12
)

func main() {

	wg.Add(2)
	go ProcessHandle("./test")
	go NotifyConstant()

	wg.Wait()
}

func NotifyConstant() {
	defer wg.Done()
	for {
		if err := Notify(con_text); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Hour * 12)
	}
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

	text := fmt.Sprintf("%s proccessが止まりました。", bin)
	if err := Notify(text); err != nil {
		log.Fatal(err)
	}
}

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
