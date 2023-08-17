package win

import (
	"syscall"
	"unsafe"
)

const (
	WM_CREATE        = 0x0001
	WM_DESTROY       = 0x0002
	WM_PAINT         = 0x000F
	WM_SIZE          = 0x0005
	WM_DISPLAYCHANGE = 0x007E
)

const (
	CS_VREDRAW = 0x0001
	CS_HREDRAW = 0x0002
)

const (
	IDC_ARROW = 32512
)

const CW_USEDEFAULT = 0x80000000

const WS_OVERLAPPEDWINDOW = 0x00CF0000

const SW_NORMAL = 1

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

func HIWORD(dwValue uint32) uint16 {
	return uint16(dwValue >> 16)
}

func LOWORD(dwValue uint32) uint16 {
	return uint16(dwValue)
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterClassEx  = user32.NewProc("RegisterClassExW")
	procDefWindowProc    = user32.NewProc("DefWindowProcW")
	procLoadCursor       = user32.NewProc("LoadCursorW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procCreateWindowEx   = user32.NewProc("CreateWindowExW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procUpdateWindow     = user32.NewProc("UpdateWindow")
	procGetClientRect    = user32.NewProc("GetClientRect")
	procValidateRect     = user32.NewProc("ValidateRect")
	procInvalidateRect   = user32.NewProc("InvalidateRect")
)

func RegisterClassEx(wndClassEx *WNDCLASSEX) (uint16, error) {
	r, _, err := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wndClassEx)))
	if r == 0 {
		return 0, err
	}
	return uint16(r), nil
}

func DefWindowProc(hWnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	r, _, _ := procDefWindowProc.Call(uintptr(hWnd), uintptr(msg), wparam, lparam)
	return r
}

func LoadCursor(hInstance syscall.Handle, lpCursorName uint32) (syscall.Handle, error) {
	r, _, err := procLoadCursor.Call(uintptr(hInstance), uintptr(lpCursorName))
	if r == 0 {
		return syscall.InvalidHandle, err
	}
	return syscall.Handle(r), nil
}

func PostQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode))
}

func GetMessage(msg *MSG, hWnd syscall.Handle, msgFilterMin, msgFilterMax uint32) bool {
	r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(msg)), uintptr(hWnd), uintptr(msgFilterMin), uintptr(msgFilterMax))
	return r != 0
}

func TranslateMessage(lpMsg *MSG) bool {
	r, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(lpMsg)))
	return r != 0
}

func DispatchMessage(lpMsg *MSG) (uintptr, error) {
	r, _, err := procDispatchMessage.Call(uintptr(unsafe.Pointer(lpMsg)))
	if r == 0 {
		return 0, err
	}
	return r, nil
}

func CreateWindowEx(dwExStyle uint32, lpClassName, lpWindowName string, dwStyle uint32, x, y, nWidth, nHeight uint32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (syscall.Handle, error) {
	lpClassName16, err := syscall.UTF16PtrFromString(lpClassName)
	if err != nil {
		return syscall.InvalidHandle, err
	}
	lpWindowName16, err := syscall.UTF16PtrFromString(lpWindowName)
	if err != nil {
		return syscall.InvalidHandle, err
	}
	r, _, err := procCreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(lpClassName16)),
		uintptr(unsafe.Pointer(lpWindowName16)),
		uintptr(dwStyle),
		uintptr(x),
		uintptr(y),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam),
	)
	if r == 0 {
		return syscall.InvalidHandle, err
	}
	return syscall.Handle(r), nil
}

func ShowWindow(hWnd syscall.Handle, cmdShow int32) bool {
	r, _, _ := procShowWindow.Call(uintptr(hWnd), uintptr(cmdShow))
	return r != 0
}

func UpdateWindow(hWnd syscall.Handle) bool {
	r, _, _ := procUpdateWindow.Call(uintptr(hWnd))
	return r != 0
}

func GetClientRect(hWnd syscall.Handle, lpRect *RECT) bool {
	r, _, _ := procGetClientRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(lpRect)))
	return r != 0
}

func ValidateRect(hWnd syscall.Handle, lpRect *RECT) bool {
	r, _, _ := procValidateRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(lpRect)))
	return r != 0
}

func InvalidateRect(hWnd syscall.Handle, lpRect *RECT, bErase bool) bool {
	erase := 0
	if bErase {
		erase = 1
	}
	r, _, _ := procInvalidateRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(lpRect)), uintptr(erase))
	return r != 0
}
