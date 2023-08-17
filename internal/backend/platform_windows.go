package backend

import (
	"runtime"
	"syscall"
	"unsafe"

	"gni.dev/gni"
	"gni.dev/gni/internal/win"
)

const className = "gni.dev app"

var rootWidget gni.Widget
var d2d *Direct2D

func Exec(w gni.Widget) {
	runtime.LockOSThread()

	rootWidget = w

	var err error

	d2d, err = NewDirect2D()
	if err != nil {
		panic("d2d initialization error: " + err.Error())
	}

	hInst, err := win.GetModuleHandle()
	if err != nil {
		panic("GetModuleHandle failed" + err.Error())
	}

	hCursor, err := win.LoadCursor(0, win.IDC_ARROW)
	if err != nil {
		panic("LoadCursor failed: " + err.Error())
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hInst
	wc.HCursor = hCursor
	wc.LpszClassName, _ = syscall.UTF16PtrFromString(className)
	if _, err := win.RegisterClassEx(&wc); err != nil {
		panic("RegisterClassEx failed: " + err.Error())
	}

	hWnd, err := win.CreateWindowEx(0, className, "GNI", win.WS_OVERLAPPEDWINDOW, win.CW_USEDEFAULT, win.CW_USEDEFAULT, 800, 600, 0, 0, hInst, 0)
	if err != nil {
		panic("CreateWindow failed: " + err.Error())
	}

	win.ShowWindow(hWnd, win.SW_NORMAL)
	win.UpdateWindow(hWnd)

	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	d2d.Release()
}

func wndProc(hWnd syscall.Handle, uMsg uint32, wParam, lParam uintptr) uintptr {
	switch uMsg {
	case win.WM_DISPLAYCHANGE:
		win.InvalidateRect(hWnd, nil, false)
		return 0
	case win.WM_PAINT:
		onRender(hWnd)
		win.ValidateRect(hWnd, nil)
		return 0
	case win.WM_SIZE:
		width := win.LOWORD(uint32(lParam))
		height := win.HIWORD(uint32(lParam))
		d2d.Resize(width, height)
		return 0
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 1
	default:
		return win.DefWindowProc(hWnd, uMsg, wParam, lParam)
	}
}

func onRender(hWnd syscall.Handle) {
	if err := d2d.BeginDraw(hWnd); err != nil {
		return
	}
	rootWidget.OnPaint(d2d)
	d2d.EndDraw()
}
