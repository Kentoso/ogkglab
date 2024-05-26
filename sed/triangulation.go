package sed

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"strconv"
)

type Triangle struct {
	vertices  [3]Point
	neighbors []int // Indices of adjacent triangles in the TriangulationResult
}

type TriangulationResult struct {
	originalVertices []Point
	triangles        []Triangle
}

type Corridor struct {
	vertices []Point
	doors    [][2]Point // doors of the corridor
}

func (c *Corridor) Print() {
	fmt.Println("Corridor vertices:")
	for _, vertex := range c.vertices {
		fmt.Printf("(%.2f, %.2f) ", vertex.X, vertex.Y)
	}
	fmt.Println("\nCorridor doors:")
	for _, door := range c.doors {
		fmt.Printf("((%.2f, %.2f), (%.2f, %.2f)) ", door[0].X, door[0].Y, door[1].X, door[1].Y)
	}
	fmt.Println()
}

func (t *TriangulationResult) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	for i, triangle := range t.triangles {
		for i := 0; i < len(triangle.vertices); i++ {
			start := triangle.vertices[i]
			end := triangle.vertices[(i+1)%len(triangle.vertices)]
			line := canvas.NewLine(color.Black)
			line.Move(fyne.NewPos(start.X, start.Y))
			line.Resize(fyne.NewSize(end.X-start.X, end.Y-start.Y))
			objects = append(objects, line)
		}

		centroidX := (triangle.vertices[0].X + triangle.vertices[1].X + triangle.vertices[2].X) / 3
		centroidY := (triangle.vertices[0].Y + triangle.vertices[1].Y + triangle.vertices[2].Y) / 3

		// Create and position the text with the triangle index
		text := canvas.NewText(strconv.Itoa(i), color.Black)
		text.Move(fyne.NewPos(centroidX-5., centroidY-5.))
		objects = append(objects, text)
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
	var triangles []Triangle
	var triangleIndices [][]int

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
					triangles = append(triangles, Triangle{vertices: [3]Point{prev, curr, next}})
					triangleIndices = append(triangleIndices, []int{(i + n - 1) % n, i, (i + 1) % n})
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
		triangles = append(triangles, Triangle{vertices: [3]Point{vertices[0], vertices[1], vertices[2]}})
		triangleIndices = append(triangleIndices, []int{0, 1, 2})
	}

	return &TriangulationResult{originalVertices: polygon.vertices, triangles: triangles}
}

func assignNeighbors(triangles []Triangle) {
	for i := 0; i < len(triangles); i++ {
		for j := i + 1; j < len(triangles); j++ {
			sharedVertices := 0
			for a := 0; a < 3; a++ {
				for b := 0; b < 3; b++ {
					if pointsEqual(triangles[i].vertices[a], triangles[j].vertices[b]) {
						sharedVertices++
					}
				}
			}
			if sharedVertices == 2 {
				triangles[i].neighbors = append(triangles[i].neighbors, j)
				triangles[j].neighbors = append(triangles[j].neighbors, i)
			}
		}
	}
}

func IncorporatePoints(tr *TriangulationResult, s, t Point) *TriangulationResult {
	addedPoints := []Point{s, t}
	for _, point := range addedPoints {
		for i, triangle := range tr.triangles {
			if pointInTriangle(point, triangle.vertices[0], triangle.vertices[1], triangle.vertices[2]) {
				newTriangle1 := Triangle{
					vertices:  [3]Point{point, triangle.vertices[0], triangle.vertices[1]},
					neighbors: []int{},
				}
				newTriangle2 := Triangle{
					vertices:  [3]Point{point, triangle.vertices[1], triangle.vertices[2]},
					neighbors: []int{},
				}
				newTriangle3 := Triangle{
					vertices:  [3]Point{point, triangle.vertices[2], triangle.vertices[0]},
					neighbors: []int{},
				}

				tr.triangles = append(tr.triangles, newTriangle1, newTriangle2, newTriangle3)
				tr.triangles = append(tr.triangles[:i], tr.triangles[i+1:]...)
				break
			}
		}
	}

	assignNeighbors(tr.triangles)

	return &TriangulationResult{
		originalVertices: append(tr.originalVertices, s, t),
		triangles:        tr.triangles,
	}
}

func (tr *TriangulationResult) GetCorridors() []Corridor {
	dualGraph := BuildGraphDual(*tr)

	simplifiedDualGraph := dualGraph.Simplify()

	// Identify junction triangles (degree >= 3)
	junctions := map[int]bool{}
	for node, _ := range simplifiedDualGraph.adjacencyList {
		junctions[node] = true
	}

	corridors := []Corridor{}
	visited := map[int]bool{}
	for i := 0; i < len(tr.triangles); i++ {
		if !visited[i] && !junctions[i] {
			corridor := Corridor{vertices: []Point{}, doors: [][2]Point{}}
			queue := []int{i}
			for len(queue) > 0 {
				curr := queue[0]
				queue = queue[1:]
				if visited[curr] || junctions[curr] {
					continue
				}
				visited[curr] = true
				corridor.vertices = append(corridor.vertices, tr.triangles[curr].vertices[:]...)
				for _, neighbor := range tr.triangles[curr].neighbors {
					if !visited[neighbor] && !junctions[neighbor] {
						queue = append(queue, neighbor)
					} else if junctions[neighbor] {
						corridor.doors = append(corridor.doors, [2]Point{tr.triangles[curr].vertices[0], tr.triangles[neighbor].vertices[0]})
					}
				}
			}
			corridors = append(corridors, corridor)
		}
	}
	return corridors
}
