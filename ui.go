package main

import (
	//"fmt"
	"strings"
	"termtyper/key"

	"github.com/zserge/webview"
)

var w webview.WebView

func searchandpaste(datapath string) {
	w = webview.New(settings.Termtyper.Debug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Bind("snipAsyncRequest", snipAsyncRequest)
	w.Bind("snipFromClip", snipGetClipboard)
	w.Bind("snipTyper", snipTyper)
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)

	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	if settings.Termtyper.WebLocal == false {
		w.Navigate("https://nimblearchitect.github.io/termtyper/common/searchpage.html")
	} else {
		w.Navigate("file://" + datapath + "/common/searchpage.html")
	}
	//w.Init(snipSearchRemote())
	w.Run()
}

func newfromcommand(datapath string) {
	w = webview.New(settings.Termtyper.Debug)
	defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)
	w.Bind("snipCodeFromArg", snipCodeFromArg)

	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	w.Navigate("https://nimblearchitect.github.io/termtyper/common/searchpage.html")
	w.Eval("window.addEventListener('load', function () { getCodeFromArguments(); });")
	//w.Navigate("file://" + datapath + "/common/createnew.html")
	w.Run()
}

func typemanager(datapath string) {
	w = webview.New(settings.Termtyper.Debug)
	//defer w.Destroy()
	w.SetTitle(appName)
	w.SetSize(1024, 768, webview.HintNone)
	w.Bind("snipAsyncRequest", snipAsyncRequest)
	w.Bind("snipFromClip", snipCopy)
	w.Bind("snipClose", snipClose)
	w.Bind("snipSave", snipSave)

	w.Navigate("data:text/html,<html><body>Loading...</body></html>")
	//w.Navigate("https://nimblearchitect.github.io/termtyper/common/manager.html")
	w.Navigate("file://" + datapath + "/common/manager.html")
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

func minimizeWindow() {
	ptr := w.Window()
	key.SwitchWindow(ptr)
}
