package sedv2

import "math"

type Point struct {
	X, Y float32
}

func (p Point) Distance(q Point) float32 {
	return float32(math.Sqrt(float64((p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y))))
}

func (p Point) Angle(q Point) float32 {
	dx := q.X - p.X
	dy := q.Y - p.Y
	angle := math.Atan2(float64(dy), float64(dx))
	if angle < 0 {
		angle += 2 * math.Pi // Ensure the angle is in the range [0, 2π)
	}
	return float32(angle)
}

func (p Point) OnSegment(a, b Point) bool {
	if p.X <= max(a.X, b.X) && p.X >= min(a.X, b.X) &&
		p.Y <= max(a.Y, b.Y) && p.Y >= min(a.Y, b.Y) {
		crossProduct := (p.Y-a.Y)*(b.X-a.X) - (p.X-a.X)*(b.Y-a.Y)
		if crossProduct == 0 {
			return true
		}
	}
	return false
}

func (p Point) GetIntersectionWithRay(raydir Point, q Point, s Point) (Point, bool) {
	rx, ry := raydir.X, raydir.Y
	sx, sy := s.X-q.X, s.Y-q.Y
	qpx, qpy := q.X-p.X, q.Y-p.Y

	det := rx*sy - ry*sx
	if math.Abs(float64(det)) < 1e-10 {
		return Point{}, false
	}

	t := (qpx*sy - qpy*sx) / det
	u := (qpx*ry - qpy*rx) / det

	if t >= 0 && u >= 0 && u <= 1 {
		ix := p.X + t*rx
		iy := p.Y + t*ry
		return Point{ix, iy}, true
	}

	return Point{}, false
}
