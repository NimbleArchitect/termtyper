package key

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "sendkeys_windows.h"
import "C"
import (
	"time"
	"unsafe"
)

func SwitchWindow() {

}

func SendLine(text string, keyDelay int) {
	delay := time.Duration(keyDelay)

	for _, c := range text {
		code := string(c)

		//mod := C.int(shift)
		name := C.CString(string(code))
		defer C.free(unsafe.Pointer(name))

		C.Sendkey(name, 0)
		time.Sleep(delay * time.Millisecond)
	}
}
