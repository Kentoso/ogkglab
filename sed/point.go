package sed

type Point struct {
	X, Y float32
}

func (p Point) Distance(q Point) float32 {
	return (p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y)
}
