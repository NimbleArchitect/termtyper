package main

import (
	//"fmt"
	"github.com/zserge/webview"
	"strings"
)

var w webview.WebView

func searchandpaste(datapath string) {
	w = webview.New(webdebug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	//w.Navigate("https://nimblearchitect.github.io/termtyper/common/searchpage.html")
	w.Navigate("file://" + datapath + "/common/searchpage.html")
	w.Bind("snipAsyncSearch", snipAsyncSearch)
	w.Bind("snipFromClip", snipGetClipboard)
	w.Bind("snipWrite", snipWrite)
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)

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
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)
	w.Bind("snipCodeFromArg", snipCodeFromArg)
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
	w.Bind("snipAsyncSearch", snipAsyncSearch)
	w.Bind("snipFromClip", snipCopy)
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)

	w.Run()
}

func sendResultsToJS(hash string, results string) {
	//need to escape the esacpes and escape the single quotes
	out := strings.Replace(results, "\\", "\\\\", -1)
	out = strings.Replace(out, "'", "\\'", -1)

	w.Dispatch(func() {
		w.Eval("asyncJob.GotData('" + hash + "','" + out + "');")
	})
}
