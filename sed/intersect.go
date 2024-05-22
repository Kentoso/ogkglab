package sed

func isBetween(a, b, c Point) bool {
	return (a.X <= c.X && c.X <= b.X || b.X <= c.X && c.X <= a.X) &&
		(a.Y <= c.Y && c.Y <= b.Y || b.Y <= c.Y && c.Y <= a.Y)
}

func crossProduct(a, b, c Point) float32 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func segmentsIntersect(p1, q1, p2, q2 Point) bool {
	// Calculate the four cross products
	d1 := crossProduct(p2, q2, p1)
	d2 := crossProduct(p2, q2, q1)
	d3 := crossProduct(p1, q1, p2)
	d4 := crossProduct(p1, q1, q2)

	// Check if the line segments intersect
	if d1*d2 < 0 && d3*d4 < 0 {
		return true
	}

	// Special case: Check if the intersection occurs at endpoints
	if d1 == 0 && isBetween(p2, q2, p1) {
		return false
	}
	if d2 == 0 && isBetween(p2, q2, q1) {
		return false
	}
	if d3 == 0 && isBetween(p1, q1, p2) {
		return false
	}
	if d4 == 0 && isBetween(p1, q1, q2) {
		return false
	}

	return false
}
