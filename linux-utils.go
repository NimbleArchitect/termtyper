// +build linux

package main

import (
	"github.com/go-vgo/robotgo"
	"time"
)

func typeSnippet(text string) {
	w.Minimize()
	time.Sleep(1 * time.Second)
	robotgo.TypeStr(text)
	time.Sleep(1 * time.Second)
	w.Terminate()
}
