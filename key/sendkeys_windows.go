package key

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -L. -lXtst
// #include <stdlib.h>
// #include "sendkeys_windows.h"
import "C"
import (
	"time"
	"unsafe"
)

func SwitchWindow() {

}

func SendLine(text string) {

	for _, c := range text {
		code, shift := char2keyCode(string(c))

		mod := C.int(shift)
		name := C.CString(string(code))
		defer C.free(unsafe.Pointer(name))

		C.Sendkey(name, mod)
	}
}
