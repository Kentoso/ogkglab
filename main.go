package main

import (
	"fmt"
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
	polygonMap   *sedv2.Map
	window       *fyne.Window
	menu         *fyne.Container
	currentState int
}

type DrawingObject interface {
	Draw() fyne.CanvasObject
}

func drawObject(o DrawingObject) *fyne.Container {
	return container.NewCenter(o.Draw())
}

func getMenu(nextButtonFunc func()) *fyne.Container {
	obstaclesInput := widget.NewMultiLineEntry()
	obstaclesInput.Resize(fyne.NewSize(50, 50))
	sInput := widget.NewEntry()
	tInput := widget.NewEntry()
	startAndTargetInput := container.NewVBox(sInput, tInput)
	nextButton := widget.NewButton("Next", nextButtonFunc)
	menu := container.NewGridWithColumns(3, container.NewHBox(nextButton), obstaclesInput, startAndTargetInput)
	return menu
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
	updateWindow(game, drawObject(polygonMap))
	polygonMap.FindShortestPath()
	fmt.Println(polygonMap.Results.Path)
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
	//polygonMap.AddObstacles(sedv2.Obstacle{Vertices: []sedv2.Point{{10, 10}, {10, 90}, {90, 90}, {90, 10}}})
	polygonMap.AddObstacles(sedv2.CreateRandomObstacle(3, -90, -90, 0, 0))
	polygonMap.AddObstacles(sedv2.CreateRandomObstacle(3, 10, 10, 80, 80))

	g := Game{
		polygonMap:   polygonMap,
		window:       &window,
		currentState: stateInput,
	}

	nextStateFunc := func() {
		g.currentState = (g.currentState + 1) % len(stateFuncs)
		stateFuncs[g.currentState](&g, polygonMap)
	}

	g.menu = getMenu(nextStateFunc)

	stateFuncs[g.currentState](&g, polygonMap)

	fmt.Println(window.Content())

	window.Resize(fyne.NewSize(500, 500))
	window.ShowAndRun()
}
