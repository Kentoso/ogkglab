package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"ogkglab/sedv2"
)

const (
	stateInput = iota
	stateMap
	stateVisibilityGraph
	stateShortestPath
)

//var currentState = stateInput

type Game struct {
	polygonMap     *sedv2.Map
	window         *fyne.Window
	menu           *fyne.Container
	currentState   int
	sInput         *widget.Entry
	tInput         *widget.Entry
	obstaclesInput *widget.Entry
}

type DrawingObject interface {
	Draw() fyne.CanvasObject
}

func drawObject(o DrawingObject) *InteractiveCanvas {
	return NewInteractiveCanvas(o.Draw())
}

func getMenu(nextButtonFunc func(), randomButtonFunc func()) (*fyne.Container, *widget.Entry, *widget.Entry, *widget.Entry) {
	obstaclesInput := widget.NewMultiLineEntry()
	obstaclesInput.Resize(fyne.NewSize(50, 50))
	sInput := widget.NewEntry()
	tInput := widget.NewEntry()
	startAndTargetInput := container.NewVBox(sInput, tInput)
	nextButton := widget.NewButton("Next", nextButtonFunc)
	randomButton := widget.NewButton("Random", randomButtonFunc)
	menu := container.NewGridWithColumns(3, container.NewHBox(nextButton, randomButton), obstaclesInput, startAndTargetInput)
	return menu, sInput, tInput, obstaclesInput
}

func updateWindow(game *Game, drawing *InteractiveCanvas) {
	menu := game.menu
	window := game.window

	var drawingContainer fyne.CanvasObject

	background := canvas.NewRectangle(color.White)
	if drawing == nil {
		drawingContainer = container.NewStack(background)
	} else {
		drawingContainer = container.NewStack(background, drawing)
	}

	content := container.NewBorder(menu, nil, nil, nil, drawingContainer)

	(*window).SetContent(content)
}

func inputState(game *Game, _ *sedv2.Map) {
	updateWindow(game, NewInteractiveCanvas(nil))
}

func mapState(game *Game, polygonMap *sedv2.Map) {
	if game.obstaclesInput.Text != "" && game.sInput.Text != "" && game.tInput.Text != "" {
		polygonMap.Clear()
		polygonMap.AddObstacles(parseObstacles(game.obstaclesInput.Text)...)
		polygonMap.S, _ = parsePoint(game.sInput.Text)
		polygonMap.T, _ = parsePoint(game.tInput.Text)
	}
	updateWindow(game, drawObject(polygonMap))
	polygonMap.FindShortestPath()
}

func visibilityGraphState(game *Game, polygonMap *sedv2.Map) {
	//updateWindow(game, drawObject(polygonMap.StepsResults.Polygon))
	updateWindow(game, drawObject(polygonMap.Results.VisibilityGraph))
}

func shortestPathState(game *Game, polygonMap *sedv2.Map) {
	updateWindow(game, drawObject(polygonMap))
}

var stateFuncs = []func(*Game, *sedv2.Map){inputState, mapState, visibilityGraphState, shortestPathState}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Border Layout")

	polygonMap := sedv2.NewMap(sedv2.Point{100, 100}, sedv2.Point{200, 200})

	obstacles := []sedv2.Obstacle{
		//sedv2.Obstacle{Vertices: []sedv2.Point{{-}},
		sedv2.CreateRandomObstacle(3, -90, -90, -30, -30).Translate(200, 200),
		sedv2.CreateRandomObstacle(99, 40, 40, 80, 80).Translate(200, 200),
		//sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}},
		//sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}}.Translate(-90, -90),
		//sedv2.Obstacle{Vertices: []sedv2.Point{{-80, -70}, {-40, -30}, {-60, -60}}},
		//sedv2.Obstacle{Vertices: []sedv2.Point{{40, 0}, {60, 70}, {30, 60}}},
	}

	polygonMap.AddObstacles(obstacles...)

	g := Game{
		polygonMap:   polygonMap,
		window:       &window,
		currentState: stateInput,
	}

	nextStateFunc := func() {
		g.currentState = (g.currentState + 1) % len(stateFuncs)
		stateFuncs[g.currentState](&g, polygonMap)
	}

	randomButtonFunc := func() {
		if g.currentState == stateInput {
			g.obstaclesInput.SetText("")
			S, _ := parsePoint(g.sInput.Text)
			T, _ := parsePoint(g.tInput.Text)
			var minX, minY, maxX, maxY float32 = math.MaxFloat32, math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32
			for _, p := range []sedv2.Point{S, T} {
				minX = float32(math.Min(float64(minX), float64(p.X)))
				minY = float32(math.Min(float64(minY), float64(p.Y)))
				maxX = float32(math.Max(float64(maxX), float64(p.X)))
				maxY = float32(math.Max(float64(maxY), float64(p.Y)))
			}

			obstacleNumX, obstacleNumY := int(math.Sqrt(float64(maxX-minX))), int(math.Sqrt(float64(maxY-minY)))
			sizeX, sizeY := int((maxX-minX)/float32(obstacleNumX)), int((maxY-minY)/float32(obstacleNumY))
			maxObstacles := obstacleNumX * obstacleNumY
			currentObstacles := 0
			for i := 0; i < obstacleNumX; i++ {
				if currentObstacles >= maxObstacles {
					break
				}
				for j := 0; j < obstacleNumY; j++ {
					oMinX, oMinY, oMaxX, oMaxY := minX+float32(i*sizeX), minY+float32(j*sizeY), minX+float32((i+1)*sizeX), minY+float32((j+1)*sizeY)

					currentObstacles++
					g.obstaclesInput.SetText(g.obstaclesInput.Text + "\n\n" + sedv2.CreateRandomObstacle(3, oMinX, oMinY, oMaxX, oMaxY).ToString())
					if currentObstacles >= maxObstacles {
						break
					}
				}
			}
		}
	}

	g.menu, g.sInput, g.tInput, g.obstaclesInput = getMenu(nextStateFunc, randomButtonFunc)

	g.sInput.Text, g.tInput.Text = "100,100", "200,200"

	stateFuncs[g.currentState](&g, polygonMap)

	window.Resize(fyne.NewSize(500, 500))
	window.ShowAndRun()
}
