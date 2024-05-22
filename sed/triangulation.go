package sed

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
)

type Triangle struct {
	vertices [3]Point
}

type TriangulationResult struct {
	originalVertices []Point
	edges            []Segment
	triangles        []Triangle
}

func (t *TriangulationResult) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	for _, edge := range t.edges {
		line := canvas.NewLine(color.Black)
		line.Move(fyne.NewPos(edge.start.X, edge.start.Y))
		line.Resize(fyne.NewSize(edge.end.X-edge.start.X, edge.end.Y-edge.start.Y))
		objects = append(objects, line)
	}

	for _, vertex := range t.originalVertices {
		circle := canvas.NewCircle(color.Black)
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(vertex.X-2.5, vertex.Y-2.5))
		objects = append(objects, circle)
	}

	return container.NewWithoutLayout(objects...)
}

func angleBetween(p1, p2, p3 Point) float64 {
	a := math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X))
	b := math.Atan2(float64(p3.Y-p2.Y), float64(p3.X-p2.X))
	angle := b - a
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func isEar(p *Polygon, indices []int, i int) bool {
	n := len(indices)
	if n < 3 {
		return false
	}
	prev := (i + n - 1) % n
	next := (i + 1) % n
	tri := []Point{p.vertices[indices[prev]], p.vertices[indices[i]], p.vertices[indices[next]]}
	for j := 0; j < n; j++ {
		if j == prev || j == i || j == next {
			continue
		}
		if pointInTriangle(p.vertices[indices[j]], tri) {
			return false
		}
	}
	return true
}

func pointInTriangle(p Point, tri []Point) bool {
	d1 := crossProduct(tri[0], tri[1], p)
	d2 := crossProduct(tri[1], tri[2], p)
	d3 := crossProduct(tri[2], tri[0], p)
	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	return !(hasNeg && hasPos)
}

func TriangulateEarClipping(p Polygon) *TriangulationResult {
	n := len(p.vertices)
	if n < 3 {
		return nil
	}

	originalVertices := make([]Point, len(p.vertices))
	copy(originalVertices, p.vertices)

	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}

	// Step 1: Compute the interior angles of each vertex
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		prev := (i + n - 1) % n
		next := (i + 1) % n
		angles[i] = angleBetween(p.vertices[indices[prev]], p.vertices[indices[i]], p.vertices[indices[next]])
	}

	// Step 2: Identify each vertex as an ear tip or not
	isEarTip := make([]bool, n)
	for i := 0; i < n; i++ {
		isEarTip[i] = isEar(&p, indices, i)
	}

	triangles := []Triangle{}
	edges := []Segment{}

	// Step 3: Find the ear tip with the smallest interior angle and construct triangles
	for len(triangles) < len(originalVertices)-2 {
		minAngle := math.Inf(1)
		earTip := -1
		for i := 0; i < n; i++ {
			if isEarTip[i] && angles[i] < minAngle {
				minAngle = angles[i]
				earTip = i
			}
		}

		if earTip == -1 {
			break
		}

		prev := (earTip + n - 1) % n
		next := (earTip + 1) % n
		triangle := Triangle{vertices: [3]Point{p.vertices[indices[prev]], p.vertices[indices[earTip]], p.vertices[indices[next]]}}
		triangles = append(triangles, triangle)
		edges = append(edges, Segment{p.vertices[indices[prev]], p.vertices[indices[earTip]]})
		edges = append(edges, Segment{p.vertices[indices[earTip]], p.vertices[indices[next]]})
		edges = append(edges, Segment{p.vertices[indices[next]], p.vertices[indices[prev]]})

		// Remove the ear tip vertex
		indices = append(indices[:earTip], indices[earTip+1:]...)
		n--

		// Update the angles and ear tip status of the affected vertices
		if n > 0 {
			prev = (earTip + n - 1) % n
			next = earTip % n
			angles[prev] = angleBetween(p.vertices[indices[(prev+n-1)%n]], p.vertices[indices[prev]], p.vertices[indices[(prev+1)%n]])
			isEarTip[prev] = isEar(&p, indices, prev)
			angles[next] = angleBetween(p.vertices[indices[(next+n-1)%n]], p.vertices[indices[next]], p.vertices[indices[(next+1)%n]])
			isEarTip[next] = isEar(&p, indices, next)
		}
	}

	return &TriangulationResult{
		originalVertices: originalVertices,
		edges:            edges,
		triangles:        triangles,
	}
}
