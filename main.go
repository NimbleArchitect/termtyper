// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/go-vgo/robotgo
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	//"fmt"
	"github.com/atotto/clipboard"
	"github.com/zserge/webview"
)

const debug = true

var w webview.WebView
var action int

func main() {

	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	searchandpaste()

}

func searchandpaste() {
	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("snip search")
	w.SetSize(600, 400, webview.HintNone)
	//w.Navigate("data:text/html," + html)
	w.Navigate("file:///home/rich/data/src/go/src/snippets/frontpage.html")
	w.Bind("searchsnip", searchsnip)
	w.Bind("toclipboard", copysnip)
	w.Bind("writesnip", writesnip)
	w.Bind("closesnip", closesnip)
	w.Run()
}

func searchsnip(data string) string {
	//time.Sleep(4 * time.Second)
	//println("running from js: " + data)
	return `[{"id": 1456, "txt": "mike"}, {"id": 25672, "txt":"dave"}]`
}

func copysnip(data string) error {
	clipboard.WriteAll(data)
	return nil
}

func writesnip(data string) error {
	go typeSnippet(data)
	return nil
}

func closesnip() error {
	go w.Terminate()
	return nil
}
