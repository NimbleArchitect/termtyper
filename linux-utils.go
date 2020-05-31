// +build linux

package main

import (
	"github.com/go-vgo/robotgo"
	"time"
)

func typeSnippet(text []string) {
	w.Minimize()
	time.Sleep(1 * time.Second)

	count := len(text)
	for i := 0; i < count; i++ {
		//fmt.Println(scanner.Text())
		singleline := text[i]
		robotgo.TypeStr(singleline)
		if i < (count - 1) {
			robotgo.KeyTap("enter")
			time.Sleep(100 * time.Millisecond)
		}
	}
	time.Sleep(1 * time.Second)
	w.Terminate()
}
