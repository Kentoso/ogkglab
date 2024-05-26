package sed

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

type TriangulationResult struct {
	originalVertices []Point
	triangles        []Polygon
}

func (t *TriangulationResult) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	for _, triangle := range t.triangles {
		for i := 0; i < len(triangle.vertices); i++ {
			start := triangle.vertices[i]
			end := triangle.vertices[(i+1)%len(triangle.vertices)]
			line := canvas.NewLine(color.Black)
			line.Move(fyne.NewPos(start.X, start.Y))
			line.Resize(fyne.NewSize(end.X-start.X, end.Y-start.Y))
			objects = append(objects, line)
		}
	}

	for _, vertex := range t.originalVertices {
		circle := canvas.NewCircle(color.Black)
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(vertex.X-2.5, vertex.Y-2.5))
		objects = append(objects, circle)
	}

	return container.NewWithoutLayout(objects...)
}

func pointInTriangle(pt, v1, v2, v3 Point) bool {
	// Barycentric coordinate method
	d1 := sign(pt, v1, v2)
	d2 := sign(pt, v2, v3)
	d3 := sign(pt, v3, v1)
	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 Point) float32 {
	return (p2.X-p1.X)*(p3.Y-p1.Y) - (p3.X-p1.X)*(p2.Y-p1.Y)
}

func isConvex(a, b, c Point) bool {
	return sign(a, b, c) < 0
}

func pointsEqual(p1, p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

func TriangulateEarClipping(polygon Polygon) *TriangulationResult {
	if len(polygon.vertices) < 3 {
		return &TriangulationResult{originalVertices: polygon.vertices, triangles: nil} // No triangulation possible
	}

	vertices := append([]Point(nil), polygon.vertices...)
	triangles := []Polygon{}

	for len(vertices) > 3 {
		earFound := false
		n := len(vertices)

		for i := 0; i < n; i++ {
			prev := vertices[(i+n-1)%n]
			curr := vertices[i]
			next := vertices[(i+1)%n]

			if isConvex(prev, curr, next) {
				ear := true

				for j := 0; j < n; j++ {
					if j == (i+n-1)%n || j == i || j == (i+1)%n ||
						pointsEqual(vertices[j], vertices[(i+n-1)%n]) ||
						pointsEqual(vertices[j], vertices[i]) ||
						pointsEqual(vertices[j], vertices[(i+1)%n]) {
						continue
					}
					if pointInTriangle(vertices[j], prev, curr, next) {
						ear = false
						break
					}
				}

				if ear {
					triangles = append(triangles, Polygon{vertices: []Point{prev, curr, next}})
					vertices = append(vertices[:i], vertices[i+1:]...)
					earFound = true
					break
				}
			}
		}

		if !earFound {
			break // No ear found, probably a complex polygon
		}
	}

	// Add the last remaining triangle
	if len(vertices) == 3 {
		triangles = append(triangles, Polygon{vertices: vertices})
	}

	return &TriangulationResult{originalVertices: polygon.vertices, triangles: triangles}
}
