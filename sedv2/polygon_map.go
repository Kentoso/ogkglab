package sedv2

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Results struct {
	VisibilityGraph *VisibilityGraph
	Path            []Point
}

type Obstacle struct {
	Vertices      []Point
	sortingVertex Point
}

type Map struct {
	obstacles []Obstacle
	S         Point
	T         Point
	Results   Results
}

func NewMap(S, T Point) *Map {
	return &Map{nil, S, T, Results{}}
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

	if m.Results.Path != nil {
		for i := 0; i < len(m.Results.Path)-1; i++ {
			start := m.Results.Path[i]
			end := m.Results.Path[i+1]
			line := canvas.NewLine(color.RGBA{0, 255, 0, 255})
			line.Position1 = fyne.NewPos(start.X, start.Y)
			line.Position2 = fyne.NewPos(end.X, end.Y)
			objects = append(objects, line)
		}
	}

	return container.NewWithoutLayout(objects...)
}

func (m *Map) AddObstacles(obstacles ...Obstacle) {
	m.obstacles = append(m.obstacles, obstacles...)
}

func (m *Map) FindShortestPath() []Point {
	visibilityGraph := GetVisibilityGraph(m.obstacles, m.S, m.T)
	m.Results.VisibilityGraph = &visibilityGraph
	_, path := visibilityGraph.ShortestEuclideanDistance()
	m.Results.Path = path
	return path
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
