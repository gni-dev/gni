package graphics

// Point is a 2D point.
type Point struct {
	X, Y float32
}

func Pt[T number](x, y T) Point {
	return Point{float32(x), float32(y)}
}

// Rectangle is a 2D rectangle.
type Rectangle struct {
	X0, Y0, X1, Y1 float32
}

func Rect[T number](x0, y0, x1, y1 T) Rectangle {
	return Rectangle{float32(x0), float32(y0), float32(x1), float32(y1)}
}
