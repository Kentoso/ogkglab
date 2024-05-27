package main

import (
	"bufio"
	"fmt"
	"ogkglab/sedv2"
	"strconv"
	"strings"
)

func parsePoint(line string) (sedv2.Point, error) {
	coords := strings.Split(line, ",")
	if len(coords) != 2 {
		return sedv2.Point{}, fmt.Errorf("invalid format")
	}

	x, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 32)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 32)

	if err1 != nil || err2 != nil {
		return sedv2.Point{}, fmt.Errorf("invalid coordinates")
	}

	return sedv2.Point{X: float32(x), Y: float32(y)}, nil
}

func parseObstacles(input string) []sedv2.Obstacle {
	var obstacles []sedv2.Obstacle
	var vertices []sedv2.Point

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(vertices) > 0 {
				obstacles = append(obstacles, sedv2.Obstacle{Vertices: vertices})
				vertices = []sedv2.Point{}
			}
			continue
		}

		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			continue
		}

		x, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 32)
		y, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 32)

		if err1 == nil && err2 == nil {
			vertices = append(vertices, sedv2.Point{X: float32(x), Y: float32(y)})
		}
	}

	if len(vertices) > 0 {
		obstacles = append(obstacles, sedv2.Obstacle{Vertices: vertices})
	}

	return obstacles
}
