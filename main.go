package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"ogkglab/sed"
)

const (
	stateInput = iota
	stateMap
	statePolygon
	stateTriangulation
)

//var currentState = stateInput

type Game struct {
	polygonMap   *sed.Map
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
	//menu := getMenu(nextButtonFunc)
	menu := game.menu
	window := game.window

	background := canvas.NewRectangle(color.White)
	drawingContainer := container.NewStack(background, drawing)
	content := container.NewBorder(menu, nil, nil, nil, drawingContainer)

	(*window).SetContent(content)
}

func inputState(game *Game, polygonMap *sed.Map) {
	updateWindow(game, drawObject(polygonMap))
}

func mapState(game *Game, polygonMap *sed.Map) {
	updateWindow(game, drawObject(polygonMap))
}

func polygonState(game *Game, polygonMap *sed.Map) {
	updateWindow(game, drawObject(polygonMap.ToPolygon()))
}

func triangulationState(game *Game, polygonMap *sed.Map) {
	updateWindow(game, drawObject(sed.TriangulateEarClipping(*polygonMap.ToPolygon())))

}

var stateFuncs = []func(*Game, *sed.Map){inputState, mapState, polygonState, triangulationState}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Border Layout")

	polygonMap := sed.NewMap(sed.Point{10, 10}, sed.Point{100, 100})
	polygonMap.AddObstacle(sed.CreateRandomObstacle(3, 10, 10, 100, 100))

	g := Game{
		polygonMap:   polygonMap,
		window:       &window,
		currentState: statePolygon,
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
