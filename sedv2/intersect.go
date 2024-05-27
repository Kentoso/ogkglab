package sedv2

func isBetween(a, b, c Point) bool {
	return (a.X <= c.X && c.X <= b.X || b.X <= c.X && c.X <= a.X) &&
		(a.Y <= c.Y && c.Y <= b.Y || b.Y <= c.Y && c.Y <= a.Y)
}

func crossProduct(a, b, c Point) float32 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func doSegmentsIntersect(p1, q1, p2, q2 Point) bool {
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
	if d1 == 0 && p2 != p1 && q2 != p1 && isBetween(p2, q2, p1) {
		return true
	}
	if d2 == 0 && p2 != q1 && q2 != q1 && isBetween(p2, q2, q1) {
		return true
	}
	if d3 == 0 && p1 != p2 && q1 != p2 && isBetween(p1, q1, p2) {
		return true
	}
	if d4 == 0 && p1 != q2 && q1 != q2 && isBetween(p1, q1, q2) {
		return true
	}

	return false
}

func doSegmentsIntersectAlternative(p1, q1, p2, q2 Point) bool {
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
		return true
	}
	if d2 == 0 && isBetween(p2, q2, q1) {
		return true
	}
	if d3 == 0 && isBetween(p1, q1, p2) {
		return true
	}
	if d4 == 0 && isBetween(p1, q1, q2) {
		return true
	}

	return false
}

func doesSegmentIntersectObstacle(a, b Point, obstacle Obstacle) bool {
	for i := 0; i < len(obstacle.Vertices); i++ {
		start := obstacle.Vertices[i]
		end := obstacle.Vertices[(i+1)%len(obstacle.Vertices)]
		if doSegmentsIntersect(a, b, start, end) {
			return true
		}
	}

	return false
}
