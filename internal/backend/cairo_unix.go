//go:build (linux && !android) || freebsd || openbsd

package backend

/*
#cgo pkg-config: cairo

#include <cairo.h>
#include <stdlib.h>
*/
import "C"
import (
	"math"
	"unsafe"

	"gni.dev/gni/graphics"
)

type Cairo struct {
	cr *C.cairo_t
}

func (c *Cairo) SetAntiAlias(aa bool) {
	if aa {
		C.cairo_set_antialias(c.cr, C.CAIRO_ANTIALIAS_DEFAULT)
	} else {
		C.cairo_set_antialias(c.cr, C.CAIRO_ANTIALIAS_NONE)
	}
}

func (c *Cairo) DrawRoundRect(rect graphics.Rectangle, rx, ry float32) {
	C.cairo_new_sub_path(c.cr)
	c.arcRxRy(rect.X1-rx, rect.Y0+ry, rx, ry, deg2rad(-90), deg2rad(0))
	c.arcRxRy(rect.X1-rx, rect.Y1-ry, rx, ry, deg2rad(0), deg2rad(90))
	c.arcRxRy(rect.X0+rx, rect.Y1-ry, rx, ry, deg2rad(90), deg2rad(180))
	c.arcRxRy(rect.X0+rx, rect.Y0+ry, rx, ry, deg2rad(180), deg2rad(270))
	C.cairo_close_path(c.cr)
	C.cairo_stroke(c.cr)
}

func (c *Cairo) DrawText(text string, at graphics.Point) {
	C.cairo_move_to(c.cr, C.double(at.X), C.double(at.Y))

	cs := C.CString(text)
	C.cairo_show_text(c.cr, cs)
	C.free(unsafe.Pointer(cs))
}

func (c *Cairo) SetFontSize(size float32) {
	C.cairo_set_font_size(c.cr, C.double(size))
}

func (c *Cairo) arcRxRy(cx, cy, rx, ry, rad1, rad2 float32) {
	if rx == ry {
		C.cairo_arc(c.cr, C.double(cx), C.double(cy), C.double(rx), C.double(rad1), C.double(rad2))
	} else {
		C.cairo_save(c.cr)
		C.cairo_translate(c.cr, C.double(cx), C.double(cy))
		C.cairo_scale(c.cr, C.double(rx), C.double(ry))
		C.cairo_arc(c.cr, 0, 0, 1, C.double(rad1), C.double(rad2))
		C.cairo_restore(c.cr)
	}
}

func deg2rad(deg float32) float32 {
	return deg * math.Pi / 180.0
}
