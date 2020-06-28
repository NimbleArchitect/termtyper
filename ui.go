package main

import (
	"github.com/zserge/webview"
)

var w webview.WebView

func searchandpaste(datapath string) {
	w = webview.New(webdebug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintMin)
	//w.Navigate("data:text/html," + html)
	w.Navigate("https://nimblearchitect.github.io/termtyper/common/searchpage.html")
	//w.Navigate("file://" + datapath + "/searchpage.html")
	w.Bind("snipSearch", snip_search)
	w.Bind("toclipboard", snip_copy)
	w.Bind("snipWrite", snip_write)
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)
	w.Run()
}

func newfromcommand(datapath string) {
	w = webview.New(webdebug)
	defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://nimblearchitect.github.io/termtyper/common/createnew.html")
	w.Navigate("file://" + datapath + "/createnew.html")
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)
	w.Bind("snipCodeFromArg", snip_codeFromArg)
	w.Eval("window.addEventListener('load', function () { getCodeFromArguments(); });")
	w.Run()
}
