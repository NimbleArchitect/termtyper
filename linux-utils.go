// +build linux

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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

	//fire separate thread so we can send to stdin
	go func() {
		//send keys to type to stdin of python script :(
		count := len(text)
		for i := 0; i < count; i++ {
			singleline := text[i]
			if i < (count - 1) { //more than one line and we are not on the last
				stdin.Write([]byte(singleline + " \\\n")) //sent line of text
				stdin.Write([]byte("{ENTER}\n"))          //now move to the next line
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

func readStdin() string {
	var retstr string

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("No Pipe found")
		//return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []string

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			output = append(output, string(input))
			break
		}
		output = append(output, string(input))
	}

	for j := 0; j < len(output); j++ {
		retstr += output[j]
	}
	return strings.TrimSpace(retstr)
}
