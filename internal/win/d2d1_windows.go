package win

import (
	"math"
	"syscall"
	"unsafe"
)

const S_OK = 0

const D2DERR_RECREATE_TARGET = 0x8899000C

type D2D1_FACTORY_TYPE int32

const D2D1_FACTORY_TYPE_SINGLE_THREADED D2D1_FACTORY_TYPE = 0

type D2D1_RENDER_TARGET_TYPE int32

const D2D1_RENDER_TARGET_TYPE_DEFAULT D2D1_RENDER_TARGET_TYPE = 0

type DXGI_FORMAT int32

const DXGI_FORMAT_UNKNOWN DXGI_FORMAT = 0

type D2D1_ALPHA_MODE int32

const D2D1_ALPHA_MODE_UNKNOWN D2D1_ALPHA_MODE = 0

type D2D1_RENDER_TARGET_USAGE int32

const D2D1_RENDER_TARGET_USAGE_NONE D2D1_RENDER_TARGET_USAGE = 0

type D2D1_FEATURE_LEVEL int32

const D2D1_FEATURE_LEVEL_DEFAULT D2D1_FEATURE_LEVEL = 0

type D2D1_PRESENT_OPTIONS int32

const D2D1_PRESENT_OPTIONS_NONE D2D1_PRESENT_OPTIONS = 0

type D2D1_ANTIALIAS_MODE int32

const (
	D2D1_ANTIALIAS_MODE_PER_PRIMITIVE D2D1_ANTIALIAS_MODE = 0
	D2D1_ANTIALIAS_MODE_ALIASED       D2D1_ANTIALIAS_MODE = 1
)

type D2D1_TEXT_ANTIALIAS_MODE int32

const (
	D2D1_TEXT_ANTIALIAS_MODE_DEFAULT D2D1_TEXT_ANTIALIAS_MODE = 0
	D2D1_TEXT_ANTIALIAS_MODE_ALIASED D2D1_TEXT_ANTIALIAS_MODE = 3
)

type D2D1_PIXEL_FORMAT struct {
	Format    DXGI_FORMAT
	AlphaMode D2D1_ALPHA_MODE
}

type D2D1_RENDER_TARGET_PROPERTIES struct {
	Type        D2D1_RENDER_TARGET_TYPE
	PixelFormat D2D1_PIXEL_FORMAT
	DpiX        float32
	DpiY        float32
	Usage       D2D1_RENDER_TARGET_USAGE
	MinLevel    D2D1_FEATURE_LEVEL
}

type D2D1_SIZE_U struct {
	Width  uint32
	Height uint32
}

type D2D1_HWND_RENDER_TARGET_PROPERTIES struct {
	Hwnd           syscall.Handle
	PixelSize      D2D1_SIZE_U
	PresentOptions D2D1_PRESENT_OPTIONS
}

type D2D1_POINT_2F struct {
	X, Y float32
}

type D2D1_RECT_F struct {
	Left, Top, Right, Bottom float32
}

type D2D1_ROUNDED_RECT struct {
	Rect             D2D1_RECT_F
	RadiusX, RadiusY float32
}

type D2D1_COLOR_F struct {
	R, G, B, A float32
}

type D2D1_MATRIX_3X2_F struct {
	A1, B1, A2, B2, A3, B3 float32
}

type D2D1_BRUSH_PROPERTIES struct {
	Opacity   float32
	Transform D2D1_MATRIX_3X2_F
}

var (
	d2d1                  = syscall.NewLazyDLL("d2d1.dll")
	procD2D1CreateFactory = d2d1.NewProc("D2D1CreateFactory")
)

var (
	IID_ID2D1Factory = syscall.GUID{Data1: 0x06152247, Data2: 0x6f50, Data3: 0x465a, Data4: [8]byte{0x92, 0x45, 0x11, 0x8b, 0xfd, 0x3b, 0x60, 0x07}}
)

type ID2D1Factory struct {
	vtbl *ID2D1FactoryVtbl
}

type ID2D1FactoryVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	ReloadSystemMetrics            uintptr
	GetDesktopDpi                  uintptr
	CreateRectangleGeometry        uintptr
	CreateRoundedRectangleGeometry uintptr
	CreateEllipseGeometry          uintptr
	CreateGeometryGroup            uintptr
	CreateTransformedGeometry      uintptr
	CreatePathGeometry             uintptr
	CreateStrokeStyle              uintptr
	CreateDrawingStateBlock        uintptr
	CreateWicBitmapRenderTarget    uintptr
	CreateHwndRenderTarget         uintptr
	CreateDxgiSurfaceRenderTarget  uintptr
	CreateDCRenderTarget           uintptr
}

func D2D1CreateFactory(factoryType D2D1_FACTORY_TYPE, riid *syscall.GUID) (*ID2D1Factory, error) {
	var pFactory *ID2D1Factory
	r, _, err := procD2D1CreateFactory.Call(
		uintptr(factoryType),
		uintptr(unsafe.Pointer(riid)),
		0,
		uintptr(unsafe.Pointer(&pFactory)),
	)
	if r != S_OK {
		return nil, err
	}
	return pFactory, nil
}

func (i *ID2D1Factory) Release() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *ID2D1Factory) CreateHwndRenderTarget(renderTargetProperties *D2D1_RENDER_TARGET_PROPERTIES, hwndRenderTargetProperties *D2D1_HWND_RENDER_TARGET_PROPERTIES) (*ID2D1HwndRenderTarget, error) {
	var hwndRenderTarget *ID2D1HwndRenderTarget
	r, _, err := syscall.SyscallN(i.vtbl.CreateHwndRenderTarget, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(renderTargetProperties)),
		uintptr(unsafe.Pointer(hwndRenderTargetProperties)),
		uintptr(unsafe.Pointer(&hwndRenderTarget)),
	)
	if r != S_OK {
		return nil, err
	}
	return hwndRenderTarget, nil
}

type ID2D1HwndRenderTarget struct {
	vtbl *ID2D1HwndRenderTargetVtbl
}

type ID2D1HwndRenderTargetVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetFactory uintptr

	CreateBitmap                 uintptr
	CreateBitmapFromWicBitmap    uintptr
	CreateSharedBitmap           uintptr
	CreateBitmapBrush            uintptr
	CreateSolidColorBrush        uintptr
	CreateGradientStopCollection uintptr
	CreateLinearGradientBrush    uintptr
	CreateRadialGradientBrush    uintptr
	CreateCompatibleRenderTarget uintptr
	CreateLayer                  uintptr
	CreateMesh                   uintptr
	DrawLine                     uintptr
	DrawRectangle                uintptr
	FillRectangle                uintptr
	DrawRoundedRectangle         uintptr
	FillRoundedRectangle         uintptr
	DrawEllipse                  uintptr
	FillEllipse                  uintptr
	DrawGeometry                 uintptr
	FillGeometry                 uintptr
	FillMesh                     uintptr
	FillOpacityMask              uintptr
	DrawBitmap                   uintptr
	DrawText                     uintptr
	DrawTextLayout               uintptr
	DrawGlyphRun                 uintptr
	SetTransform                 uintptr
	GetTransform                 uintptr
	SetAntialiasMode             uintptr
	GetAntialiasMode             uintptr
	SetTextAntialiasMode         uintptr
	GetTextAntialiasMode         uintptr
	SetTextRenderingParams       uintptr
	GetTextRenderingParams       uintptr
	SetTags                      uintptr
	GetTags                      uintptr
	PushLayer                    uintptr
	PopLayer                     uintptr
	Flush                        uintptr
	SaveDrawingState             uintptr
	RestoreDrawingState          uintptr
	PushAxisAlignedClip          uintptr
	PopAxisAlignedClip           uintptr
	Clear                        uintptr
	BeginDraw                    uintptr
	EndDraw                      uintptr
	GetPixelFormat               uintptr
	SetDpi                       uintptr
	GetDpi                       uintptr
	GetSize                      uintptr
	GetPixelSize                 uintptr
	GetMaximumBitmapSize         uintptr
	IsSupported                  uintptr

	CheckWindowState uintptr
	Resize           uintptr
	GetHwnd          uintptr
}

func (i *ID2D1HwndRenderTarget) Release() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *ID2D1HwndRenderTarget) CreateSolidColorBrush(color *D2D1_COLOR_F, brushProperties *D2D1_BRUSH_PROPERTIES) (*ID2D1SolidColorBrush, error) {
	var solidColorBrush *ID2D1SolidColorBrush
	r, _, err := syscall.SyscallN(i.vtbl.CreateSolidColorBrush, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(color)),
		uintptr(unsafe.Pointer(brushProperties)),
		uintptr(unsafe.Pointer(&solidColorBrush)),
	)
	if r != S_OK {
		return nil, err
	}
	return solidColorBrush, nil
}

func (i *ID2D1HwndRenderTarget) DrawRoundedRectangle(rect *D2D1_ROUNDED_RECT, brush unsafe.Pointer, strokeWidth float32) {
	syscall.SyscallN(i.vtbl.DrawRoundedRectangle, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(rect)),
		uintptr(brush),
		uintptr(math.Float32bits(strokeWidth)),
		0,
	)
}

func (i *ID2D1HwndRenderTarget) DrawText(text string, textFormat *IDWriteTextFormat, layoutRect *D2D1_RECT_F, brush unsafe.Pointer) {
	text16, _ := syscall.UTF16PtrFromString(text)
	syscall.SyscallN(i.vtbl.DrawText, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(text16)),
		uintptr(len(text)),
		uintptr(unsafe.Pointer(textFormat)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(brush)),
		0,
		0,
	)
}

func (i *ID2D1HwndRenderTarget) SetTransform(transform *D2D1_MATRIX_3X2_F) {
	syscall.SyscallN(i.vtbl.SetTransform, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(transform)))
}

func (i *ID2D1HwndRenderTarget) SetAntiAliasMode(antiAliasMode D2D1_ANTIALIAS_MODE) {
	syscall.SyscallN(i.vtbl.SetAntialiasMode, uintptr(unsafe.Pointer(i)), uintptr(antiAliasMode))
}

func (i *ID2D1HwndRenderTarget) SetTextAntialiasMode(textAntiAliasMode D2D1_TEXT_ANTIALIAS_MODE) {
	syscall.SyscallN(i.vtbl.SetTextAntialiasMode, uintptr(unsafe.Pointer(i)), uintptr(textAntiAliasMode))
}

func (i *ID2D1HwndRenderTarget) Clear(color *D2D1_COLOR_F) {
	syscall.SyscallN(i.vtbl.Clear, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(color)))
}

func (i *ID2D1HwndRenderTarget) BeginDraw() {
	syscall.SyscallN(i.vtbl.BeginDraw, uintptr(unsafe.Pointer(i)))
}

func (i *ID2D1HwndRenderTarget) EndDraw() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.EndDraw, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *ID2D1HwndRenderTarget) Resize(size *D2D1_SIZE_U) {
	syscall.SyscallN(i.vtbl.Resize, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(size)))
}

type ID2D1Brush struct {
	vtbl *ID2D1BrushVtbl
}

type ID2D1BrushVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetFactory uintptr

	SetOpacity   uintptr
	SetTransform uintptr
	GetOpacity   uintptr
	GetTransform uintptr
}

type ID2D1SolidColorBrush struct {
	vtbl *ID2D1SolidColorBrushVtbl
}

type ID2D1SolidColorBrushVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetFactory uintptr

	SetOpacity   uintptr
	SetTransform uintptr
	GetOpacity   uintptr
	GetTransform uintptr

	SetColor uintptr
	GetColor uintptr
}

func (i *ID2D1SolidColorBrush) Release() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}
