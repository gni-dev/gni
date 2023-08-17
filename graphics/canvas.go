package graphics

// Canvas is the interface that wraps the basic drawing methods.
type Canvas interface {
	DrawRoundRect(rect Rectangle, rx, ry float32)
	DrawText(text string, at Point)
}
