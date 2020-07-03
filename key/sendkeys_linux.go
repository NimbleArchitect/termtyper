package key

// #cgo linux openbsd freebsd pkg-config: x11
// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -L. -lXtst
// #include <stdlib.h>
// #include "sendkeys_linux.h"
import "C"
import (
	//"fmt"
	"unsafe"
)

func SwitchWindow() {
	C.SendAltTabKeys()

}

func SendLine(text string) {
	//fmt.Println("F:SendLine:start")
	for _, c := range text {
		code, shift := char2keyCode(string(c))

		mod := C.int(shift)
		name := C.CString(string(code))
		defer C.free(unsafe.Pointer(name))
		//fmt.Println("F:SendLine:sendkey =", c)

		C.Sendkey(name, mod)
	}
	//fmt.Println("F:SendLine:end")

}

func char2keyCode(charCode string) (string, int) {
	var mod int = 0
	var key string = ""

	if len(charCode) <= 0 {
		return "", -1
	}
	switch charCode {
	//key names from https://www.cl.cam.ac.uk/~mgk25/ucs/keysymdef.h
	case " ":
		key = "space"
	case "!":
		key = "exclam"
		mod = 1
	case "\"":
		key = "quotedbl"
		mod = 1
	case "#":
		key = "numbersign"
		mod = 0
	case "$":
		key = "dollar"
		mod = 1
	case "%":
		key = "percent"
		mod = 1
	case "&":
		key = "ampersand"
		mod = 1
	case "'":
		key = "apostrophe"
	case "(":
		key = "parenleft"
		mod = 1
	case ")":
		key = "parenright"
		mod = 1
	case "*":
		key = "asterisk"
		mod = 1
	case "+":
		key = "plus"
		mod = 1
	case ",":
		key = "comma"
	case "-":
		key = "minus"
	case ".":
		key = "period"
	case "/":
		key = "slash"
	case ":":
		key = "colon"
		mod = 1
	case ";":
		key = "semicolon"
	case "<":
		key = "less"
		mod = 1
	case "=":
		key = "equal"
	case ">":
		key = "greater"
		mod = 1
	case "?":
		key = "question"
		mod = 1
	case "@":
		key = "at"
		mod = 1
	case "A":
		key = "a"
		mod = 1
	case "B":
		key = "b"
		mod = 1
	case "C":
		key = "c"
		mod = 1
	case "D":
		key = "d"
		mod = 1
	case "E":
		key = "e"
		mod = 1
	case "F":
		key = "f"
		mod = 1
	case "G":
		key = "g"
		mod = 1
	case "H":
		key = "h"
		mod = 1
	case "I":
		key = "i"
		mod = 1
	case "J":
		key = "j"
		mod = 1
	case "K":
		key = "k"
		mod = 1
	case "L":
		key = "l"
		mod = 1
	case "M":
		key = "m"
		mod = 1
	case "N":
		key = "n"
		mod = 1
	case "O":
		key = "o"
		mod = 1
	case "P":
		key = "p"
		mod = 1
	case "Q":
		key = "q"
		mod = 1
	case "R":
		key = "r"
		mod = 1
	case "S":
		key = "s"
		mod = 1
	case "T":
		key = "t"
		mod = 1
	case "U":
		key = "u"
		mod = 1
	case "V":
		key = "v"
		mod = 1
	case "W":
		key = "w"
		mod = 1
	case "X":
		key = "x"
		mod = 1
	case "Y":
		key = "y"
		mod = 1
	case "Z":
		key = "z"
		mod = 1
	case "[":
		key = "bracketleft"
	case "\\":
		key = "backslash"
	case "]":
		key = "bracketright"
	case "^":
		key = "asciicircum"
		mod = 1
	case "_":
		key = "underscore"
		mod = 1
	case "`":
		key = "grave"
	case "{":
		key = "braceleft"
		mod = 1
	case "|":
		key = "bar"
		mod = 1
	case "}":
		key = "braceright"
		mod = 1
	case "~":
		key = "asciitilde"
		mod = 1
	case "£":
		key = "sterling"
		mod = 1
	case "\n":
		key = "Return"

	default:
		key = charCode
	}

	return key, mod
}
