// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func typeSnippet(text []string) {

	//start hacky python program
	execpath := getprogPath()
	cmd := exec.Command(execpath + "/key.py")
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if nil != err {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	if err := cmd.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}

	timer := time.AfterFunc(6*time.Second, func() {
		cmd.Process.Kill()
	})

	//fire seperate thread so we can send to stdin
	go func() {
		//send keys to type to stdin of python script :(
		count := len(text)
		for i := 0; i < count; i++ {
			singleline := text[i]
			if i < (count - 1) { //more than one line and we are not on the last
				stdin.Write([]byte(singleline + "\n")) //sent line of text
				stdin.Write([]byte("{ENTER}\n"))       //now move to the next line
			} else {
				stdin.Write([]byte(singleline + "\n")) //write the last or only line
			}
		}
		//now we sent the exit string, so python quits it's loop
		stdin.Write([]byte("{TIME2QUIT}\n"))
	}()

	err = cmd.Wait()
	timer.Stop()

	w.Terminate()
}
