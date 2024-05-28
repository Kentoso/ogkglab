package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
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

func drawObject(o DrawingObject) *fyne.Container {
	return container.NewCenter(o.Draw())
}

func getMenu(nextButtonFunc func()) (*fyne.Container, *widget.Entry, *widget.Entry, *widget.Entry) {
	obstaclesInput := widget.NewMultiLineEntry()
	obstaclesInput.Resize(fyne.NewSize(50, 50))
	sInput := widget.NewEntry()
	tInput := widget.NewEntry()
	startAndTargetInput := container.NewVBox(sInput, tInput)
	nextButton := widget.NewButton("Next", nextButtonFunc)
	menu := container.NewGridWithColumns(3, container.NewHBox(nextButton), obstaclesInput, startAndTargetInput)
	return menu, sInput, tInput, obstaclesInput
}

func updateWindow(game *Game, drawing *fyne.Container) {
	menu := game.menu
	window := game.window

	background := canvas.NewRectangle(color.White)
	drawingContainer := container.NewStack(background, drawing)
	content := container.NewBorder(menu, nil, nil, nil, drawingContainer)

	(*window).SetContent(content)
}

func inputState(game *Game, _ *sedv2.Map) {
	updateWindow(game, container.NewWithoutLayout())
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

	polygonMap := sedv2.NewMap(sedv2.Point{-100, -100}, sedv2.Point{100, 100})
	//polygonMap.AddObstacles(sedv2.CreateRandomObstacle(5, -90, -90, 0, 0))

	//polygonMap.AddObstacles()

	//polygonMap.AddObstacles(sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}})
	obstacles := []sedv2.Obstacle{
		//sedv2.Obstacle{Vertices: []sedv2.Point{{-}},
		sedv2.CreateRandomObstacle(3, -90, -90, -30, -30),
		sedv2.CreateRandomObstacle(3, 40, 40, 80, 80),
		//sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}},
		//sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}}.Translate(-90, -90),
		//sedv2.Obstacle{Vertices: []sedv2.Point{{-80, -70}, {-40, -30}, {-60, -60}}},
		//sedv2.Obstacle{Vertices: []sedv2.Point{{40, 0}, {60, 70}, {30, 60}}},
	}

	polygonMap.AddObstacles(obstacles...)

	//polygonMap.AddObstacles(sedv2.CreateRandomObstacle(3, -90, -90, 0, 0))
	//polygonMap.AddObstacles(sedv2.CreateRandomObstacle(3, 10, 10, 80, 80))

	g := Game{
		polygonMap:   polygonMap,
		window:       &window,
		currentState: stateInput,
	}

	nextStateFunc := func() {
		g.currentState = (g.currentState + 1) % len(stateFuncs)
		stateFuncs[g.currentState](&g, polygonMap)
	}

	g.menu, g.sInput, g.tInput, g.obstaclesInput = getMenu(nextStateFunc)

	stateFuncs[g.currentState](&g, polygonMap)

	window.Resize(fyne.NewSize(500, 500))
	window.ShowAndRun()
}
