package backend

import (
	"syscall/js"

	"gni.dev/gni/graphics"
)

type CanvasJs struct {
	this js.Value
	ctx  js.Value
}

func NewCanvasJs(doc js.Value) *CanvasJs {
	cnv := doc.Call("createElement", "canvas")
	style := cnv.Get("style")
	style.Set("position", "absolute")

	doc.Get("body").Call("appendChild", cnv)

	ctx := cnv.Call("getContext", "2d")
	return &CanvasJs{this: cnv, ctx: ctx}
}

func (c *CanvasJs) DrawRoundRect(rect graphics.Rectangle, rx, ry float32) {
	c.ctx.Call("beginPath")
	c.ctx.Call("moveTo", rect.X0+rx, rect.Y0)
	c.ctx.Call("lineTo", rect.X1-rx, rect.Y0)
	c.ctx.Call("quadraticCurveTo", rect.X1, rect.Y0, rect.X1, rect.Y0+ry)
	c.ctx.Call("lineTo", rect.X1, rect.Y1-ry)
	c.ctx.Call("quadraticCurveTo", rect.X1, rect.Y1, rect.X1-rx, rect.Y1)
	c.ctx.Call("lineTo", rect.X0+rx, rect.Y1)
	c.ctx.Call("quadraticCurveTo", rect.X0, rect.Y1, rect.X0, rect.Y1-ry)
	c.ctx.Call("lineTo", rect.X0, rect.Y0+ry)
	c.ctx.Call("quadraticCurveTo", rect.X0, rect.Y0, rect.X0+rx, rect.Y0)
	c.ctx.Call("closePath")
	c.ctx.Call("stroke")
}

func (c *CanvasJs) DrawText(text string, at graphics.Point) {
	c.ctx.Call("fillText", text, at.X, at.Y)
}

func (c *CanvasJs) Clear() {
	c.ctx.Set("lineWidth", 1)
	c.ctx.Call("clearRect", 0, 0, c.this.Get("width"), c.this.Get("height"))
	c.ctx.Call("beginPath")
}
