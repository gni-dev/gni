package backend

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"gni.dev/gni"
)

/*
#include <platform_android.h>
#include <stdlib.h>
*/
import "C"

type jvalue uint64

var (
	waitDebugger string

	rootWidget gni.Widget
)

var java struct {
	canvas struct {
		cls C.jclass

		drawRoundRect C.jmethodID
		drawText      C.jmethodID
	}

	paint struct {
		cls C.jclass

		style struct {
			cls C.jclass

			fill          C.jfieldID
			stroke        C.jfieldID
			fillAndStroke C.jfieldID
		}

		ctor           C.jmethodID
		setAntiAlias   C.jmethodID
		setStrokeWidth C.jmethodID
		setStyle       C.jmethodID
	}
}

//go:linkname callMain main.main
func callMain()

func Exec(w gni.Widget) {
	rootWidget = w
}

//export Java_dev_gni_GniActivity_gniCreate
func Java_dev_gni_GniActivity_gniCreate(env *C.JNIEnv, class C.jclass, context C.jobject) {
	if waitDebugger == "true" {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGCONT)
		<-sig
	}
	initJava(env)
	callMain()
}

//export Java_dev_gni_GniView_gniDraw
func Java_dev_gni_GniView_gniDraw(env *C.JNIEnv, class C.jclass, canvas C.jobject) {
	c := NewCanvas(env, canvas)
	rootWidget.OnPaint(c)
}

func initJava(env *C.JNIEnv) {
	java.canvas.cls = C.jclass(C.JNI_NewGlobalRef(env, C.jobject(findClass(env, "android/graphics/Canvas"))))
	java.canvas.drawRoundRect = getMethodID(env, java.canvas.cls, "drawRoundRect", "(FFFFFFLandroid/graphics/Paint;)V")
	java.canvas.drawText = getMethodID(env, java.canvas.cls, "drawText", "(Ljava/lang/String;FFLandroid/graphics/Paint;)V")

	java.paint.cls = C.jclass(C.JNI_NewGlobalRef(env, C.jobject(findClass(env, "android/graphics/Paint"))))
	java.paint.style.cls = C.jclass(C.JNI_NewGlobalRef(env, C.jobject(findClass(env, "android/graphics/Paint$Style"))))
	java.paint.style.fill = getStaticFieldID(env, java.paint.style.cls, "FILL", "Landroid/graphics/Paint$Style;")
	java.paint.style.stroke = getStaticFieldID(env, java.paint.style.cls, "STROKE", "Landroid/graphics/Paint$Style;")
	java.paint.style.fillAndStroke = getStaticFieldID(env, java.paint.style.cls, "FILL_AND_STROKE", "Landroid/graphics/Paint$Style;")
	java.paint.ctor = getMethodID(env, java.paint.cls, "<init>", "()V")
	java.paint.setAntiAlias = getMethodID(env, java.paint.cls, "setAntiAlias", "(Z)V")
	java.paint.setStrokeWidth = getMethodID(env, java.paint.cls, "setStrokeWidth", "(F)V")
	java.paint.setStyle = getMethodID(env, java.paint.cls, "setStyle", "(Landroid/graphics/Paint$Style;)V")
}

func vaArgs(args []jvalue) *C.jvalue {
	if len(args) == 0 {
		return nil
	}
	return (*C.jvalue)(unsafe.Pointer(&args[0]))
}

func jfloat(v float32) jvalue {
	return jvalue(math.Float32bits(v))
}

func jboolean(v bool) jvalue {
	if v {
		return C.JNI_TRUE
	} else {
		return C.JNI_FALSE
	}
}

func jstring(env *C.JNIEnv, v string) jvalue {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))
	return jvalue(C.JNI_NewStringUTF(env, cstr))
}

func goString(env *C.JNIEnv, str C.jstring) string {
	utf := C.JNI_GetStringUTFChars(env, str, nil)
	defer C.JNI_ReleaseStringUTFChars(env, str, utf)
	return C.GoString(utf)
}

func catchExc(env *C.JNIEnv) error {
	e := C.JNI_ExceptionOccurred(env)
	if e == 0 {
		return nil
	}
	C.JNI_ExceptionClear(env)
	cls := C.JNI_GetObjectClass(env, C.jobject(e))
	toString := getMethodID(env, cls, "toString", "()Ljava/lang/String;")
	msg := C.JNI_CallObjectMethodA(env, C.jobject(e), toString, nil)
	if C.JNI_ExceptionCheck(env) == C.JNI_TRUE {
		panic("CallObjectMethod()")
	}
	return errors.New(goString(env, C.jstring(msg)))
}

func findClass(env *C.JNIEnv, name string) C.jclass {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cls := C.JNI_FindClass(env, cName)
	if cls == 0 {
		panic(fmt.Sprintf("FindClass(%s)", name))
	}
	return cls
}

func getMethodID(env *C.JNIEnv, clazz C.jclass, name, sig string) C.jmethodID {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cSig := C.CString(sig)
	defer C.free(unsafe.Pointer(cSig))
	jm := C.JNI_GetMethodID(env, clazz, cName, cSig)
	if jm == nil {
		panic(fmt.Sprintf("GetMethodID(%s, %s)", name, sig))
	}
	return jm
}

func getStaticFieldID(env *C.JNIEnv, clazz C.jclass, name, sig string) C.jfieldID {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cSig := C.CString(sig)
	defer C.free(unsafe.Pointer(cSig))
	jf := C.JNI_GetStaticFieldID(env, clazz, cName, cSig)
	if jf == nil {
		panic(fmt.Sprintf("GetStaticFieldID(%s, %s)", name, sig))
	}
	return jf
}

func callVoidMethod(env *C.JNIEnv, obj C.jobject, method C.jmethodID, args ...jvalue) error {
	C.JNI_CallVoidMethodA(env, obj, method, vaArgs(args))
	return catchExc(env)
}

func newObject(env *C.JNIEnv, cls C.jclass, method C.jmethodID, args ...jvalue) (C.jobject, error) {
	res := C.JNI_NewObjectA(env, cls, method, vaArgs(args))
	return res, catchExc(env)
}
