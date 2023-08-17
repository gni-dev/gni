package win

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

type DWRITE_FACTORY_TYPE int32

const DWRITE_FACTORY_TYPE_SHARED DWRITE_FACTORY_TYPE = 0

type DWRITE_FONT_WEIGHT int32

const DWRITE_FONT_WEIGHT_NORMAL DWRITE_FONT_WEIGHT = 400

type DWRITE_FONT_STYLE int32

const DWRITE_FONT_STYLE_NORMAL DWRITE_FONT_STYLE = 0

type DWRITE_FONT_STRETCH int32

const DWRITE_FONT_STRETCH_NORMAL DWRITE_FONT_STRETCH = 5

type DWRITE_TEXT_ALIGNMENT int32

const DWRITE_TEXT_ALIGNMENT_CENTER DWRITE_TEXT_ALIGNMENT = 2

type DWRITE_PARAGRAPH_ALIGNMENT int32

const DWRITE_PARAGRAPH_ALIGNMENT_CENTER DWRITE_PARAGRAPH_ALIGNMENT = 1

var (
	dwrite                  = syscall.NewLazyDLL("dwrite.dll")
	procDwriteCreateFactory = dwrite.NewProc("DWriteCreateFactory")
)

var (
	IID_IDWriteFactory = syscall.GUID{Data1: 0xb859ee5a, Data2: 0xd838, Data3: 0x4b5b, Data4: [8]byte{0xa2, 0xe8, 0x1a, 0xdc, 0x7d, 0x93, 0xdb, 0x48}}
)

type IDWriteFactory struct {
	vtbl *IDWriteFactoryVtbl
}

type IDWriteFactoryVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetSystemFontCollection        uintptr
	CreateCustomFontCollection     uintptr
	RegisterFontCollectionLoader   uintptr
	UnregisterFontCollectionLoader uintptr
	CreateFontFileReference        uintptr
	CreateCustomFontFileReference  uintptr
	CreateFontFace                 uintptr
	CreateRenderingParams          uintptr
	CreateMonitorRenderingParams   uintptr
	CreateCustomRenderingParams    uintptr
	RegisterFontFileLoader         uintptr
	UnregisterFontFileLoader       uintptr
	CreateTextFormat               uintptr
	CreateTypography               uintptr
	CreateGdiInterop               uintptr
	CreateTextLayout               uintptr
	CreateGdiCompatibleTextLayout  uintptr
	CreateEllipsisTrimmingSign     uintptr
	CreateTextAnalyzer             uintptr
	CreateNumberSubstitution       uintptr
	CreateGlyphRunAnalysis         uintptr
}

func DWriteCreateFactory(factoryType DWRITE_FACTORY_TYPE, riid *syscall.GUID) (*IDWriteFactory, error) {
	var factory *IDWriteFactory
	r, _, err := procDwriteCreateFactory.Call(
		uintptr(factoryType),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&factory)),
	)
	if r != S_OK {
		return nil, err
	}
	return factory, nil
}

func (i *IDWriteFactory) Release() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *IDWriteFactory) CreateTextFormat(
	fontFamilyName string,
	fontWeight DWRITE_FONT_WEIGHT,
	fontStyle DWRITE_FONT_STYLE,
	fontStretch DWRITE_FONT_STRETCH,
	fontSize float32,
	localeName string,
) (*IDWriteTextFormat, error) {
	var textFormat *IDWriteTextFormat

	fontFamilyName16, err := syscall.UTF16PtrFromString(fontFamilyName)
	if err != nil {
		return nil, err
	}
	localeName16, err := syscall.UTF16PtrFromString(localeName)
	if err != nil {
		return nil, err
	}
	r, _, _ := syscall.SyscallN(i.vtbl.CreateTextFormat, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(fontFamilyName16)),
		0,
		uintptr(fontWeight),
		uintptr(fontStyle),
		uintptr(fontStretch),
		uintptr(math.Float32bits(fontSize)),
		uintptr(unsafe.Pointer(localeName16)),
		uintptr(unsafe.Pointer(&textFormat)),
	)
	if r != S_OK {
		return nil, fmt.Errorf("DirectDraw: CreateTextFormat failed: 0x%X", r)
	}
	return textFormat, nil
}

type IDWriteTextFormat struct {
	vtbl *IDWriteTextFormatVtbl
}

type IDWriteTextFormatVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetTextAlignment        uintptr
	SetParagraphAlignment   uintptr
	SetWordWrapping         uintptr
	SetReadingDirection     uintptr
	SetFlowDirection        uintptr
	SetIncrementalTabStop   uintptr
	SetTrimming             uintptr
	SetLineSpacing          uintptr
	GetTextAlignment        uintptr
	GetParagraphAlignment   uintptr
	GetWordWrapping         uintptr
	GetReadingDirection     uintptr
	GetFlowDirection        uintptr
	GetIncrementalTabStop   uintptr
	GetTrimming             uintptr
	GetLineSpacing          uintptr
	GetFontCollection       uintptr
	GetFontFamilyNameLength uintptr
	GetFontFamilyName       uintptr
	GetFontWeight           uintptr
	GetFontStyle            uintptr
	GetFontStretch          uintptr
	GetFontSize             uintptr
	GetLocaleNameLength     uintptr
	GetLocaleName           uintptr
}

func (i *IDWriteTextFormat) Release() uint32 {
	r, _, _ := syscall.SyscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *IDWriteTextFormat) SetTextAlignment(textAlignment DWRITE_TEXT_ALIGNMENT) {
	syscall.SyscallN(i.vtbl.SetTextAlignment, uintptr(unsafe.Pointer(i)), uintptr(textAlignment))
}

func (i *IDWriteTextFormat) SetParagraphAlignment(paragraphAlignment DWRITE_PARAGRAPH_ALIGNMENT) {
	syscall.SyscallN(i.vtbl.SetParagraphAlignment, uintptr(unsafe.Pointer(i)), uintptr(paragraphAlignment))
}
