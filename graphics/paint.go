package graphics

type Paint interface {
	SetAntiAlias(aa bool)
	SetFontSize(size float32)
}
