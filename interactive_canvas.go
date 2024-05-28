package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type InteractiveCanvas struct {
	widget.BaseWidget
	x, y    float32
	content fyne.CanvasObject
	input   *widget.Entry
}

func NewInteractiveCanvas(content fyne.CanvasObject) *InteractiveCanvas {
	if content == nil {
		content = container.NewWithoutLayout()
	}
	w := &InteractiveCanvas{content: content}
	w.ExtendBaseWidget(w)
	return w
}

func (w *InteractiveCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &interactiveCanvasRenderer{
		widget:  w,
		content: w.content,
	}
}

func (w *InteractiveCanvas) MouseUp(event *desktop.MouseEvent) {
	//fmt.Printf("Mouse released at: (%f, %f)\n", event.Position.X, event.Position.Y)
}

func (w *InteractiveCanvas) MouseDown(event *desktop.MouseEvent) {
	w.x = event.Position.X
	w.y = event.Position.Y
	fmt.Printf("Mouse clicked at: (%f, %f)\n", w.x, w.y)
	w.Refresh()
}

type interactiveCanvasRenderer struct {
	widget  *InteractiveCanvas
	content fyne.CanvasObject
}

func (r *interactiveCanvasRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *interactiveCanvasRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *interactiveCanvasRenderer) Refresh() {
	r.content.Refresh()
}

func (r *interactiveCanvasRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *interactiveCanvasRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *interactiveCanvasRenderer) Destroy() {}

func (w *InteractiveCanvas) SetContent(content fyne.CanvasObject) {
	w.content = content
	w.Refresh()
}
