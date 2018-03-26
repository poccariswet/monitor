package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	wg.Add(2)
	go ProcessHandle("./test")
	go ProcessHandle("./test2")

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
}
