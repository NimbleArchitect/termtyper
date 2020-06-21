// +build linux

package main

func typeSnippet(text []string) {

	SendAltTab()

	//send keys to type to stdin of python script :(
	count := len(text)
	for i := 0; i < count; i++ {
		singleline := text[i]
		if i < (count - 1) { //more than one line and we are not on the last
			SendKeys(singleline + " \\\n") //sent line of text
		} else {
			SendKeys(singleline) //write the last or only line
		}
	}

	w.Terminate()
}
