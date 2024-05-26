package sed

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
	"math/rand"
	"slices"
	"sort"
	"time"
)

type Obstacle struct {
	Vertices []Point
}

func CreateRandomObstacle(numPoints int, minX, minY, maxX, maxY float32) Obstacle {
	rand.Seed(time.Now().UnixNano())

	points := make([]Point, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = Point{
			X: minX + rand.Float32()*(maxX-minX),
			Y: minY + rand.Float32()*(maxY-minY),
		}
	}

	var centroid Point
	for _, p := range points {
		centroid.X += p.X
		centroid.Y += p.Y
	}
	centroid.X /= float32(numPoints)
	centroid.Y /= float32(numPoints)

	sort.Slice(points, func(i, j int) bool {
		angle1 := math.Atan2(float64(points[i].Y-centroid.Y), float64(points[i].X-centroid.X))
		angle2 := math.Atan2(float64(points[j].Y-centroid.Y), float64(points[j].X-centroid.X))
		return angle1 < angle2
	})

	return Obstacle{Vertices: points}
}

type StepsResults struct {
	Polygon             *Polygon
	Triangulation       *TriangulationResult
	TriangulationWithST *TriangulationResult
}

type Map struct {
	obstacles    []Obstacle
	S            Point
	T            Point
	ShortestPath []Point
	StepsResults *StepsResults
}

func NewMap(s, t Point) *Map {
	return &Map{nil, s, t, nil, &StepsResults{}}
}

func (m *Map) DoesSTIntersectObstacles() bool {
	for _, obstacle := range m.obstacles {
		for i := 0; i < len(obstacle.Vertices); i++ {
			start := obstacle.Vertices[i]
			end := obstacle.Vertices[(i+1)%len(obstacle.Vertices)]
			if segmentsIntersect(m.S, m.T, start, end) {
				return true
			}
		}
	}

	return false
}

func (m *Map) FindShortestPath() {
	if !m.DoesSTIntersectObstacles() {
		m.ShortestPath = []Point{m.S, m.T}
		return
	}

	p := m.ToPolygon()
	m.StepsResults.Polygon = p

	t := TriangulateEarClipping(*p)
	m.StepsResults.Triangulation = t

	tWithSandT := IncorporatePoints(t, m.S, m.T)
	m.StepsResults.TriangulationWithST = tWithSandT

	visibilityGraph := NewVisibilityGraph(tWithSandT, m.S, m.T)
	_, path := visibilityGraph.ShortestEuclideanDistance()

	m.ShortestPath = path
}

func (m *Map) AddObstacles(obstacles ...Obstacle) {
	m.obstacles = append(m.obstacles, obstacles...)
}

func (m *Map) AddObstacle(obstacle Obstacle) {
	m.obstacles = append(m.obstacles, obstacle)
}

func (m *Map) ToPolygon() *Polygon {
	var vertices, verticesForBoundingBox []Point
	var edges []Segment

	for _, obstacle := range m.obstacles {
		verticesForBoundingBox = append(verticesForBoundingBox, obstacle.Vertices...)
	}
	verticesForBoundingBox = append(verticesForBoundingBox, m.S, m.T)

	//Vertices = append(Vertices, m.S, m.T)

	minX, minY := verticesForBoundingBox[0].X, verticesForBoundingBox[0].Y
	maxX, maxY := verticesForBoundingBox[0].X, verticesForBoundingBox[0].Y

	for _, v := range verticesForBoundingBox {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	var margin float32 = 10.0
	minX -= margin
	maxX += margin
	minY -= margin
	maxY += margin

	boundingBoxVertices := []Point{
		{X: minX, Y: minY},
		{X: minX, Y: maxY},
		{X: maxX, Y: maxY},
		{X: maxX, Y: minY},
	}

	vertices = append(vertices, boundingBoxVertices...)

	for i := 0; i < len(boundingBoxVertices); i++ {
		edges = append(edges, Segment{
			start: boundingBoxVertices[i],
			end:   boundingBoxVertices[(i+1)%len(boundingBoxVertices)],
		})
	}

	edges = m.addBridges(vertices[len(vertices)-4:], edges)
	// get Vertices from edges
	var newVertices []Point

	for _, e := range edges {
		newVertices = append(newVertices, e.start)
	}

	vertices = newVertices

	return &Polygon{
		vertices: vertices,
		edges:    edges,
	}
}

func (m *Map) addBridges(vertices []Point, edges []Segment) []Segment {
	// Create bridges for each hole
	currentEdges := slices.Clone(edges)
	for i := 0; i < len(m.obstacles); i++ {
		hole := m.obstacles[i].Vertices
		//polygon := m.obstacles[0].Vertices
		minDistance := math.MaxFloat64
		var bestBridge Segment

		verticesToCheck := slices.Clone(vertices)
		for j := 0; j < len(m.obstacles); j++ {
			if i == j {
				continue
			}

			verticesToCheck = append(verticesToCheck, m.obstacles[j].Vertices...)
		}

		for _, p := range verticesToCheck {
			for _, h := range hole {
				bridge := Segment{start: p, end: h}
				valid := true
				distance := math.Hypot(float64(bridge.end.X-bridge.start.X), float64(bridge.end.Y-bridge.start.Y))

				if distance == 0 {
					continue
				}
				// Check if this bridge intersects any existing edges
				for _, e := range currentEdges {
					if segmentsIntersect(bridge.start, bridge.end, e.start, e.end) {
						valid = false
						break
					}
				}

				if valid {
					if distance > 0 && distance < minDistance {
						minDistance = distance
						bestBridge = bridge
					}
				}
			}
		}

		// Add the best bridge to the list of edges
		//bestBridgeReversed := Segment{start: bestBridge.end, end: bestBridge.start}

		//edges = append(edges, bestBridge)
		//edges = append(edges, bestBridgeReversed)

		newEdges := []Segment{}
		inserted := false

		for _, e := range currentEdges {
			newEdges = append(newEdges, e)
			if !inserted && (e.end == bestBridge.start || e.end == bestBridge.end) {
				newEdges = append(newEdges, bestBridge)

				edgesShift := 0
				for k := 0; k < len(hole); k++ {
					if hole[k] == bestBridge.end {
						edgesShift = k
						break
					}
					//newEdges = append(newEdges, Segment{start: hole[k], end: hole[(k+1)%len(hole)]})
				}

				for k := 0; k < len(hole); k++ {
					newEdges = append(newEdges, Segment{start: hole[(k+edgesShift)%len(hole)], end: hole[(k+1+edgesShift)%len(hole)]})
				}

				reverseBestBridge := Segment{start: bestBridge.end, end: bestBridge.start}
				newEdges = append(newEdges, reverseBestBridge)

				inserted = true
			}
		}

		currentEdges = newEdges
	}

	return currentEdges
}

func (m *Map) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	// Draw obstacles
	for _, obstacle := range m.obstacles {
		for i := 0; i < len(obstacle.Vertices); i++ {
			start := obstacle.Vertices[i]
			end := obstacle.Vertices[(i+1)%len(obstacle.Vertices)]
			line := canvas.NewLine(color.Black)
			line.Position1 = fyne.NewPos(start.X, start.Y)
			line.Position2 = fyne.NewPos(end.X, end.Y)
			objects = append(objects, line)
		}
	}

	// Draw start and end points
	points := []struct {
		point Point
		label string
	}{
		{m.S, "S"},
		{m.T, "T"},
	}

	for _, p := range points {
		// Draw the circle
		circle := canvas.NewCircle(color.RGBA{255, 0, 0, 255})
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(p.point.X-2.5, p.point.Y-2.5))
		objects = append(objects, circle)

		// Draw the label
		text := canvas.NewText(p.label, color.Black)
		text.TextSize = 12
		text.Move(fyne.NewPos(p.point.X+5, p.point.Y-6)) // Positioning the text near the point
		objects = append(objects, text)
	}

	if m.ShortestPath != nil {
		for i := 0; i < len(m.ShortestPath)-1; i++ {
			start := m.ShortestPath[i]
			end := m.ShortestPath[i+1]
			line := canvas.NewLine(color.RGBA{0, 255, 0, 255})
			line.Position1 = fyne.NewPos(start.X, start.Y)
			line.Position2 = fyne.NewPos(end.X, end.Y)
			objects = append(objects, line)
		}
	}

	return container.NewWithoutLayout(objects...)
}
