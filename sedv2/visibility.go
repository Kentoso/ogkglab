package sedv2

import (
	"sort"
)

func GetVisibilityGraph(S []Obstacle, start, target Point) VisibilityGraph {
	allVertices := make([]Point, 0)
	for _, obstacle := range S {
		allVertices = append(allVertices, obstacle.Vertices...)
	}
	allVertices = append(allVertices, start, target)

	visibilityGraph := NewVisibilityGraph(start, target)

	for _, v := range allVertices {
		W := VisibleVertices(v, S)

		visibilityGraph.AddEdges(v, W)
	}
	
	return visibilityGraph
}

func sortVerticesByAngle(p Point, vertices []Point) []Point {
	type Vertex struct {
		Point
		Angle    float32
		Distance float32
	}

	sortedVertices := make([]Vertex, len(vertices))
	for i, v := range vertices {
		angle := p.Angle(v)
		distance := p.Distance(v)
		sortedVertices[i] = Vertex{v, angle, distance}
	}
	sort.Slice(sortedVertices, func(i, j int) bool {
		if sortedVertices[i].Angle-sortedVertices[j].Angle < 1e-10 {
			return sortedVertices[i].Distance < sortedVertices[j].Distance
		}
		return sortedVertices[i].Angle < sortedVertices[j].Angle
	})

	sortedPoints := make([]Point, len(sortedVertices))
	for i, v := range sortedVertices {
		sortedPoints[i] = v.Point
	}

	return sortedPoints
}

func makePointToObstacleMap(S []Obstacle) map[Point]*Obstacle {
	pointToObstacle := make(map[Point]*Obstacle)
	for i := range S {
		for _, vertex := range S[i].Vertices {
			pointToObstacle[vertex] = &S[i]
		}
	}
	return pointToObstacle
}

func VisibleVertices(p Point, S []Obstacle) []Point {
	pointToObstacle := makePointToObstacleMap(S)

	// Sort the obstacle vertices according to the clockwise angle
	allVertices := make([]Point, 0)
	for _, obstacle := range S {
		allVertices = append(allVertices, obstacle.Vertices...)
	}

	sortedVertices := sortVerticesByAngle(p, allVertices)

	T := NewSegmentTree(p)
	pDir := Point{1, 0}
	for _, obstacle := range S {
		for i := 0; i < len(obstacle.Vertices); i++ {
			start := obstacle.Vertices[i]
			end := obstacle.Vertices[(i+1)%len(obstacle.Vertices)]
			_, didIntersect := p.GetIntersectionWithRay(pDir, start, end)
			if didIntersect {
				T.AddSegmentIntersection(Segment{
					start: start,
					end:   end,
				})
			}
		}
	}

	var W []Point

	wIPrev := Point{}
	wasPrevVisible := false
	for i, vertex := range sortedVertices {
		wI := vertex
		T.currRaydir = Point{wI.X - p.X, wI.Y - p.Y}

		if Visible(i, p, wIPrev, wI, pointToObstacle, T, wasPrevVisible) {
			W = append(W, wI)
			wasPrevVisible = true
		} else {
			wasPrevVisible = false
		}

		obstacle := pointToObstacle[wI]

		other1, other2 := Point{}, Point{}
		for j := 0; j < len(obstacle.Vertices); j++ {
			start := obstacle.Vertices[j]
			end := obstacle.Vertices[(j+1)%len(obstacle.Vertices)]

			if end == wI {
				other2 = start
			}

			if start == wI {
				other1 = end
			}

		}

		o1CP, o2CP := crossProduct(p, wI, other1), crossProduct(p, wI, other2)
		if o1CP < 0 {
			T.AddSegmentIntersection(Segment{
				start: wI,
				end:   other1,
			})
		} else if o1CP > 0 {
			T.RemoveSegmentIntersection(SegmentIntersection{
				segment: Segment{
					start: wI,
					end:   other1,
				},
			})
		}

		if o2CP < 0 {
			T.AddSegmentIntersection(Segment{
				start: wI,
				end:   other2,
			})
		} else if o2CP > 0 {
			T.RemoveSegmentIntersection(SegmentIntersection{
				segment: Segment{
					start: wI,
					end:   other2,
				},
			})
		}

		wIPrev = wI
	}

	return W
}

func intersectsObstacle(p, wI Point, obstacle *Obstacle) bool {
	for i := 0; i < len(obstacle.Vertices); i++ {
		a := obstacle.Vertices[i]
		b := obstacle.Vertices[(i+1)%len(obstacle.Vertices)]
		if (a != wI) && (b != wI) {
			intersected := doSegmentsIntersect(p, wI, a, b)
			if intersected {
				return true
			}
		}
	}
	return false
}

func Visible(i int, p, wIPrev, wI Point, pointToObstacle map[Point]*Obstacle, T *SegmentIntersectionTree, wasPrevVisible bool) bool {
	obstacle, ok := pointToObstacle[wI]

	if ok && !intersectsObstacle(p, wI, obstacle) {
		return false
	}

	if i == 0 || !wIPrev.OnSegment(p, wI) {
		s, exists := T.GetLeftmostSegmentIntersection()
		if exists {
			intersects := doSegmentsIntersect(p, wI, s.segment.start, s.segment.end)
			if intersects {
				return false
			}
		}
		return true
	} else if !wasPrevVisible {
		return false
	} else {
		possibleIntersections := T.FindPossibleIntersections(wIPrev, wI)
		if len(possibleIntersections) > 0 {
			return false
		}
	}

	return true
}
