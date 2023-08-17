//go:build (linux && !android) || freebsd || openbsd

package backend

/*
#cgo pkg-config: xcb
#cgo CFLAGS: -Wall

#include <cairo.h>

void gni_exec();
*/
import "C"
import "gni.dev/gni"

var rootWidget gni.Widget

func Exec(w gni.Widget) {
	rootWidget = w
	C.gni_exec()
}

//export gni_paint
func gni_paint(cr *C.cairo_t) {
	c := Cairo{cr: cr}
	rootWidget.OnPaint(&c)
}
