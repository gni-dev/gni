package backend

/*
#import <Cocoa/Cocoa.h>

static inline CGContextRef currentContext() {
	return [[NSGraphicsContext currentContext] CGContext];
}
*/
import "C"
import (
	"unsafe"

	"gni.dev/gni/graphics"
)

type CGContext struct {
	ctx  C.CGContextRef
	w, h float32
}

func NewCGContext(w, h float32) *CGContext {
	return &CGContext{ctx: C.currentContext(), w: w, h: h}
}

func (c *CGContext) DrawRoundRect(rect graphics.Rectangle, rx, ry float32) {
	rect = c.translateRect(rect)
	r := C.CGRectMake(
		C.CGFloat(rect.X0), C.CGFloat(rect.Y0),
		C.CGFloat(rect.X1-rect.X0), C.CGFloat(rect.Y1-rect.Y0),
	)
	path := C.CGPathCreateWithRoundedRect(r, C.CGFloat(rx), C.CGFloat(ry), nil)
	C.CGContextAddPath(c.ctx, path)
	C.CGContextDrawPath(c.ctx, C.kCGPathStroke)

	C.CGPathRelease(path)
}

func (c *CGContext) DrawText(text string, at graphics.Point) {
	at = c.translatePoint(at)
	str := cfString(text)
	attr := C.CFAttributedStringCreate(0, str, 0)
	line := C.CTLineCreateWithAttributedString(attr)
	C.CGContextSetTextPosition(c.ctx, C.CGFloat(at.X), C.CGFloat(at.Y))
	C.CTLineDraw(line, c.ctx)

	C.CFRelease(C.CFTypeRef(line))
	C.CFRelease(C.CFTypeRef(attr))
	C.CFRelease(C.CFTypeRef(str))
}

func (c *CGContext) translateRect(rect graphics.Rectangle) graphics.Rectangle {
	return graphics.Rect(rect.X0, c.h-rect.Y0, rect.X1, c.h-rect.Y1)
}

func (c *CGContext) translatePoint(pt graphics.Point) graphics.Point {
	return graphics.Pt(pt.X, c.h-pt.Y)
}

func cfString(str string) C.CFStringRef {
	cStr := C.CString(str)
	ptr := unsafe.Pointer(cStr)
	return C.CFStringCreateWithBytesNoCopy(C.kCFAllocatorDefault,
		(*C.UInt8)(ptr), C.CFIndex(len(str)),
		C.kCFStringEncodingUTF8,
		C.FALSE,
		C.kCFAllocatorDefault)
}
