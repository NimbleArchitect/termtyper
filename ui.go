package main

import (
	"github.com/zserge/webview"
)

var w webview.WebView

func callWebFunc(funcName string, data string) {
	jscommand := funcName + "(" + data + ");"
	w.Eval(jscommand)
}

func searchandpaste(datapath string) {
	w = webview.New(webdebug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	//w.Navigate("https://nimblearchitect.github.io/termtyper/common/searchpage.html")
	w.Navigate("file://" + datapath + "/common/searchpage.html")
	w.Bind("snipSearch", snip_search)
	w.Bind("toclipboard", snip_copy)
	w.Bind("snipWrite", snip_write)
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)

	//w.Init(snipSearchRemote())

	w.Run()
}

func newfromcommand(datapath string) {
	w = webview.New(webdebug)
	defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	w.Navigate("https://nimblearchitect.github.io/termtyper/common/createnew.html")
	//w.Navigate("file://" + datapath + "/common/createnew.html")
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)
	w.Bind("snipCodeFromArg", snip_codeFromArg)
	w.Eval("window.addEventListener('load', function () { getCodeFromArguments(); });")
	w.Run()
}

func typemanager(datapath string) {
	w = webview.New(webdebug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(1024, 768, webview.HintNone)
	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	//w.Navigate("https://nimblearchitect.github.io/termtyper/common/manager.html")
	w.Navigate("file://" + datapath + "/common/manager.html")
	w.Bind("snipSearch", snip_search)
	w.Bind("toclipboard", snip_copy)
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)

	w.Run()
}
