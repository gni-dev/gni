package main

import (
	"gni.dev/gni/app"
	"gni.dev/gni/graphics"
)

type Painter struct {
}

func (p *Painter) OnPaint(c graphics.Canvas) {
	c.DrawText("Hello world!", graphics.Pt(10, 30))
	c.DrawRoundRect(graphics.Rect(5, 5, 140, 50), 5, 20)
}

func main() {
	app.Exec(&Painter{})
}
