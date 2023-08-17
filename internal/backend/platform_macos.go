//go:build !ios && darwin

package backend

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore

void gni_exec();
*/
import "C"
import (
	"runtime"

	"gni.dev/gni"
)

var rootWidget gni.Widget

func Exec(w gni.Widget) {
	runtime.LockOSThread()

	rootWidget = w
	C.gni_exec()
}

//export gni_paint
func gni_paint(w, h float32) {
	c := NewCGContext(w, h)
	rootWidget.OnPaint(c)
}
