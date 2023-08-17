package backend

import (
	"gni.dev/gni/graphics"
)

/*
#include <platform_android.h>
#include <stdlib.h>
*/
import "C"

type CanvasAndroid struct {
	env *C.JNIEnv
	obj C.jobject

	paint C.jobject
}

func NewCanvas(env *C.JNIEnv, obj C.jobject) *CanvasAndroid {
	c := &CanvasAndroid{env: env, obj: obj}

	var err error
	c.paint, err = newObject(env, java.paint.cls, java.paint.ctor)
	if err != nil {
		return nil
	}
	callVoidMethod(c.env, c.paint, java.paint.setAntiAlias, jboolean(true))
	callVoidMethod(c.env, c.paint, java.paint.setStrokeWidth, jfloat(3))
	return c
}

func (c *CanvasAndroid) DrawRoundRect(rect graphics.Rectangle, rx, ry float32) {
	c.setStyle(java.paint.style.stroke)

	callVoidMethod(c.env, c.obj, java.canvas.drawRoundRect,
		jfloat(rect.X0),
		jfloat(rect.Y0),
		jfloat(rect.X1),
		jfloat(rect.Y1),
		jfloat(rx),
		jfloat(ry),
		jvalue(c.paint),
	)
}

func (c *CanvasAndroid) DrawText(text string, at graphics.Point) {
	callVoidMethod(c.env, c.obj, java.canvas.drawText,
		jstring(c.env, text),
		jfloat(at.X),
		jfloat(at.Y),
		jvalue(c.paint),
	)
}

func (c *CanvasAndroid) setStyle(style C.jfieldID) {
	f := C.JNI_GetStaticObjectField(c.env, java.paint.style.cls, style)
	callVoidMethod(c.env, c.paint, java.paint.setStyle, jvalue(f))

}
