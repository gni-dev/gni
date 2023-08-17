package backend

import (
	"syscall"
	"unsafe"

	"gni.dev/gni/graphics"
	"gni.dev/gni/internal/win"
)

type Direct2D struct {
	d2dFactory    *win.ID2D1Factory
	dwriteFactory *win.IDWriteFactory
	renderTarget  *win.ID2D1HwndRenderTarget
	solidBrush    *win.ID2D1SolidColorBrush
	textFormat    *win.IDWriteTextFormat
}

func NewDirect2D() (*Direct2D, error) {
	d2dFactory, err := win.D2D1CreateFactory(win.D2D1_FACTORY_TYPE_SINGLE_THREADED, &win.IID_ID2D1Factory)
	if err != nil {
		return nil, err
	}
	dwriteFactory, err := win.DWriteCreateFactory(win.DWRITE_FACTORY_TYPE_SHARED, &win.IID_IDWriteFactory)
	if err != nil {
		d2dFactory.Release()
		return nil, err
	}
	textFormat, err := dwriteFactory.CreateTextFormat(
		"Verdana",
		win.DWRITE_FONT_WEIGHT_NORMAL,
		win.DWRITE_FONT_STYLE_NORMAL,
		win.DWRITE_FONT_STRETCH_NORMAL,
		12,
		"",
	)

	if err != nil {
		dwriteFactory.Release()
		d2dFactory.Release()
		return nil, err
	}
	return &Direct2D{d2dFactory: d2dFactory, dwriteFactory: dwriteFactory, textFormat: textFormat}, nil
}

func (d *Direct2D) SetAntiAlias(aa bool) {
	if aa {
		d.renderTarget.SetAntiAliasMode(win.D2D1_ANTIALIAS_MODE_PER_PRIMITIVE)
		d.renderTarget.SetTextAntialiasMode(win.D2D1_TEXT_ANTIALIAS_MODE_DEFAULT)
	} else {
		d.renderTarget.SetAntiAliasMode(win.D2D1_ANTIALIAS_MODE_ALIASED)
		d.renderTarget.SetTextAntialiasMode(win.D2D1_TEXT_ANTIALIAS_MODE_ALIASED)
	}
}

func (d *Direct2D) DrawRoundRect(rect graphics.Rectangle, rx, ry float32) {
	d.renderTarget.DrawRoundedRectangle(&win.D2D1_ROUNDED_RECT{
		Rect: win.D2D1_RECT_F{
			Left:   rect.X0,
			Top:    rect.Y0,
			Right:  rect.X1,
			Bottom: rect.Y1,
		},
		RadiusX: rx,
		RadiusY: ry,
	}, unsafe.Pointer(d.solidBrush), 1)
}

func (d *Direct2D) DrawText(text string, at graphics.Point) {
	d.renderTarget.DrawText(
		text,
		d.textFormat,
		&win.D2D1_RECT_F{
			Left:   at.X,
			Top:    at.Y,
			Right:  200,
			Bottom: 200,
		},
		unsafe.Pointer(d.solidBrush),
	)
}

func (d *Direct2D) SetFontSize(size float32) {

}

func (d *Direct2D) Release() {
	if d.solidBrush != nil {
		d.solidBrush.Release()
		d.solidBrush = nil
	}
	if d.renderTarget != nil {
		d.renderTarget.Release()
		d.renderTarget = nil
	}
	d.d2dFactory.Release()
	d.dwriteFactory.Release()
	d.textFormat.Release()
}

func (d *Direct2D) BeginDraw(hWnd syscall.Handle) error {
	if d.renderTarget == nil {
		var err error

		var rc win.RECT
		win.GetClientRect(hWnd, &rc)
		size := win.D2D1_SIZE_U{
			Width:  uint32(rc.Right - rc.Left),
			Height: uint32(rc.Bottom - rc.Top),
		}

		props := win.D2D1_RenderTargetProperties()
		hwndProps := win.D2D1_HwndRenderTargetProperties(hWnd, size)
		d.renderTarget, err = d.d2dFactory.CreateHwndRenderTarget(&props, &hwndProps)
		if err != nil {
			return err
		}

		d.solidBrush, err = d.renderTarget.CreateSolidColorBrush(&win.D2D1_COLOR_F{
			R: 0,
			G: 0,
			B: 0,
			A: 1,
		}, nil)
		if err != nil {
			return err
		}
	}

	d.renderTarget.BeginDraw()
	m := win.D2D1_Matrix3x2F_Identity()
	d.renderTarget.SetTransform(&m)
	d.renderTarget.Clear(&win.D2D1_COLOR_F{
		R: 1,
		G: 1,
		B: 1,
		A: 1,
	})
	return nil
}

func (d *Direct2D) EndDraw() {
	if d.renderTarget == nil {
		return
	}
	hr := d.renderTarget.EndDraw()
	if hr == win.D2DERR_RECREATE_TARGET {
		d.solidBrush.Release()
		d.solidBrush = nil
		d.renderTarget.Release()
		d.renderTarget = nil
	}
}

func (d *Direct2D) Resize(width, height uint16) {
	if d.renderTarget == nil {
		return
	}
	d.renderTarget.Resize(&win.D2D1_SIZE_U{
		Width:  uint32(width),
		Height: uint32(height),
	})
}
