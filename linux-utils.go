// +build linux

package main

import (
	//"bufio"
	// "fmt"
	//"github.com/go-vgo/robotgo"
	// "io"
	"log"
	"os"
	"os/exec"
	"time"
)

func typeSnippet(text []string) {
	w.Minimize()

	time.Sleep(1 * time.Second)
	count := len(text)
	for i := 0; i < count; i++ {
		sendline(text[i])
		if i < (count - 1) {
			sendline("{ENTER}")
		}
	}

	w.Terminate()
}

//TODO: convert this to a proper sendkeys command using X11
func sendline(singleline string) {
	//start hacky python program
	cmd := exec.Command("./key.py")
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}

	go func() {
		stdin.Write([]byte(singleline + "\n"))
	}()
	cmd.Wait()
}
