package sedv2

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

type VisibilityGraph struct {
	AdjacencyMap map[Point][]Point
	S            Point
	T            Point
}

func NewVisibilityGraph(s, t Point) VisibilityGraph {
	return VisibilityGraph{
		AdjacencyMap: make(map[Point][]Point),
		S:            s,
		T:            t,
	}
}

func (vg *VisibilityGraph) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	// Draw the edges of the visibility graph
	for start, neighbors := range vg.AdjacencyMap {
		for _, end := range neighbors {
			line := canvas.NewLine(color.RGBA{0, 0, 255, 255}) // Blue color for edges
			line.Position1 = fyne.NewPos(start.X, start.Y)
			line.Position2 = fyne.NewPos(end.X, end.Y)
			objects = append(objects, line)
		}
	}

	// Draw start and end points
	points := []struct {
		point Point
		label string
		color color.Color
	}{
		{vg.S, "S", color.RGBA{255, 0, 0, 255}}, // Red color for start point
		{vg.T, "T", color.RGBA{0, 255, 0, 255}}, // Green color for end point
	}

	for _, p := range points {
		// Draw the circle
		circle := canvas.NewCircle(p.color)
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(p.point.X-2.5, p.point.Y-2.5))
		objects = append(objects, circle)

		// Draw the label
		text := canvas.NewText(p.label, color.Black)
		text.TextSize = 12
		text.Move(fyne.NewPos(p.point.X+5, p.point.Y-6)) // Positioning the text near the point
		objects = append(objects, text)
	}

	return container.NewWithoutLayout(objects...)
}

func (vg *VisibilityGraph) AddEdges(from Point, to []Point) {
	vg.AdjacencyMap[from] = to
}

func (vg *VisibilityGraph) ShortestEuclideanDistance() (map[Point]float32, []Point) {
	// Initialize the distance map and predecessor map
	distanceMap := make(map[Point]float32)
	predecessorMap := make(map[Point]Point)
	for v := range vg.AdjacencyMap {
		distanceMap[v] = float32(1<<31 - 1)
	}
	distanceMap[vg.S] = 0

	// Initialize the priority queue
	pq := NewPriorityQueue()
	pq.PushPoint(vg.S, 0)

	// ShortestEuclideanDistance's algorithm
	for !pq.IsEmpty() {
		v, _ := pq.PopPoint()

		// Relax the edges
		for _, u := range vg.AdjacencyMap[v] {
			if distanceMap[v]+v.Distance(u) < distanceMap[u] {
				distanceMap[u] = distanceMap[v] + v.Distance(u)
				predecessorMap[u] = v
				pq.PushPoint(u, distanceMap[u])
			}
		}
	}

	// Reconstruct the path from S to T
	path := []Point{}
	curr := vg.T

	// Traverse the path from T to S using the predecessor map
	for curr != vg.S {
		path = append([]Point{curr}, path...)
		curr = predecessorMap[curr]
	}
	path = append([]Point{vg.S}, path...)

	return distanceMap, path
}
