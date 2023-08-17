package gni

import (
	"gni.dev/gni/graphics"
)

// Widget is an interface that represents a graphical element that can be drawn on a canvas.
type Widget interface {
	OnPaint(c graphics.Canvas)
}
