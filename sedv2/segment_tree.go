package sedv2

import (
	"github.com/emirpasic/gods/sets/treeset"
	"math"
	"slices"
)

type Segment struct {
	start, end Point
}

type SegmentIntersection struct {
	segment Segment
}

type SegmentIntersectionTree struct {
	set        *treeset.Set
	point      Point
	currRaydir Point
}

func (s SegmentIntersection) pointRayDistance(point Point, raydir Point) (float32, bool) {
	intersection, didIntersect := point.GetIntersectionWithRay(raydir, s.segment.start, s.segment.end)
	if !didIntersect {
		return math.MaxFloat32, false
	}

	return point.Distance(intersection), true
}

func NewSegmentTree(point Point) *SegmentIntersectionTree {
	tree := &SegmentIntersectionTree{point: point}

	segmentIntersectionComparator := func(a, b interface{}) int {
		aSegment := a.(SegmentIntersection)
		bSegment := b.(SegmentIntersection)

		if aSegment.segment.start == bSegment.segment.start && aSegment.segment.end == bSegment.segment.end ||
			aSegment.segment.end == bSegment.segment.start && aSegment.segment.start == bSegment.segment.end {
			return 0
		}

		aDistance, _ := aSegment.pointRayDistance(tree.point, tree.currRaydir)
		bDistance, _ := bSegment.pointRayDistance(tree.point, tree.currRaydir)

		if aDistance < bDistance {
			return -1
		}
		if aDistance > bDistance {
			return 1
		}

		pointsA := [2]Point{aSegment.segment.start, aSegment.segment.end}
		pointsB := [2]Point{bSegment.segment.start, bSegment.segment.end}

		aSamePointIndex, bSamePointIndex := -1, -1
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				if pointsA[i] == pointsB[j] {
					aSamePointIndex, bSamePointIndex = i, j
				}
			}
		}
		aNonSamePointIndex, bNonSamePointIndex := (aSamePointIndex+1)%2, (bSamePointIndex+1)%2
		angleA := tree.point.Angle(pointsA[aNonSamePointIndex])
		angleB := tree.point.Angle(pointsB[bNonSamePointIndex])

		if angleA < angleB {
			return -1
		}
		if angleA > angleB {
			return 1
		}

		return 0
	}

	tree.set = treeset.NewWith(segmentIntersectionComparator)

	return tree
}

func (s *SegmentIntersectionTree) AddSegmentIntersection(segment Segment) {
	s.set.Add(SegmentIntersection{segment})
}

func (s *SegmentIntersectionTree) RemoveSegmentIntersection(si SegmentIntersection) {
	s.set.Remove(si)
}

func (s *SegmentIntersectionTree) GetLeftmostSegmentIntersection() (SegmentIntersection, bool) {
	values := s.set.Values()
	if len(values) == 0 {
		return SegmentIntersection{}, false
	}
	segmentIntersections := make([]SegmentIntersection, len(values))
	for i, value := range values {
		segmentIntersections[i] = value.(SegmentIntersection)
	}

	segmentIntersectionComparator := func(a, b SegmentIntersection) int {
		aSegment := a
		bSegment := b

		if aSegment.segment.start == bSegment.segment.start && aSegment.segment.end == bSegment.segment.end ||
			aSegment.segment.end == bSegment.segment.start && aSegment.segment.start == bSegment.segment.end {
			return 0
		}

		aDistance, _ := aSegment.pointRayDistance(s.point, s.currRaydir)
		bDistance, _ := bSegment.pointRayDistance(s.point, s.currRaydir)

		if aDistance < bDistance {
			return -1
		}
		if aDistance > bDistance {
			return 1
		}

		pointsA := [2]Point{aSegment.segment.start, aSegment.segment.end}
		pointsB := [2]Point{bSegment.segment.start, bSegment.segment.end}

		aSamePointIndex, bSamePointIndex := -1, -1
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				if pointsA[i] == pointsB[j] {
					aSamePointIndex, bSamePointIndex = i, j
				}
			}
		}
		aNonSamePointIndex, bNonSamePointIndex := (aSamePointIndex+1)%2, (bSamePointIndex+1)%2
		angleA := s.point.Angle(pointsA[aNonSamePointIndex])
		angleB := s.point.Angle(pointsB[bNonSamePointIndex])

		if angleA < angleB {
			return -1
		}
		if angleA > angleB {
			return 1
		}

		return 0
	}

	slices.SortFunc(segmentIntersections, segmentIntersectionComparator)

	return segmentIntersections[0], true
}

func makeSegmentIntersectionFromPoint(a Point) SegmentIntersection {
	return SegmentIntersection{Segment{a, a}}
}

func (s *SegmentIntersectionTree) FindPossibleIntersections(a, b Point) []Segment {
	aSI, bSI := makeSegmentIntersectionFromPoint(a), makeSegmentIntersectionFromPoint(b)

	s.set.Add(aSI)

	aIndex, _ := s.set.Find(func(index int, value interface{}) bool {
		return value == aSI
	})

	s.RemoveSegmentIntersection(aSI)

	s.set.Add(bSI)

	bIndex, _ := s.set.Find(func(index int, value interface{}) bool {
		return value == bSI
	})

	s.RemoveSegmentIntersection(bSI)

	var intersections []Segment

	if aIndex != -1 && bIndex != -1 {
		if aIndex > bIndex {
			aIndex, bIndex = bIndex, aIndex
		}

		for i := aIndex; i <= bIndex; i++ {
			item := s.set.Values()[i].(SegmentIntersection)
			intersections = append(intersections, item.segment)
		}
	}

	return intersections
}
