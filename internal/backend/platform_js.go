package backend

import (
	"syscall/js"

	"gni.dev/gni"
)

var (
	rootWidget gni.Widget
	cnv        *CanvasJs
)

func Exec(w gni.Widget) {
	rootWidget = w

	doc := js.Global().Get("document")
	window := js.Global().Get("window")

	cnv = NewCanvasJs(doc)
	window.Call("addEventListener", "resize", js.FuncOf(redrawCanvas))

	cnv.Clear()
	rootWidget.OnPaint(cnv)
	select {}
}

func redrawCanvas(this js.Value, args []js.Value) any {
	cnv.Clear()
	rootWidget.OnPaint(cnv)
	return nil
}
