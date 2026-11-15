package main

import (
	"time"

	"testing"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
	"github.com/nsgo-dev/nsgo/nsgotest"
)

// TestCamera fits the board, selects a screen and a step of its flow from
// the points the camera hands on, fits the selection and pans.
func TestCamera(t *testing.T) {
	app := startDesktop(t, nsgotest.Size(1200, 900))
	open(app, "Camera")
	app.WaitUntil(2*time.Second, func() bool { return cameraChart.Ready() })
	app.Screenshot("camera")
	board := app.Find("id:board").Frame.Size
	// The board shows whole: a point of it is where the camera puts it.
	at := func(id string, x, y core.Points) driver.Point {
		f := app.FindMounted("id:board-" + id).Frame
		return cameraShown.ToView(core.Point{X: f.Origin.X + x, Y: f.Origin.Y + y})
	}
	if far := cameraShown.ToView(core.Point{X: board.Width, Y: board.Height}); far.X > app.Find("id:camera").Frame.Size.Width {
		t.Fatalf("the board's far corner shows at %v", far)
	}
	app.PressAt("id:camera", at("menu", 120, 200))
	if got := app.Find("id:camera-selected").Text; got != "menu" {
		t.Fatalf("pressing the menu selected %q", got)
	}
	zoom := cameraShown.Zoom()
	app.Press("label:Fit the selection")
	app.Settle()
	if cameraShown.Zoom() <= zoom {
		t.Fatalf("fitting the selection zoomed to %v from %v", cameraShown.Zoom(), zoom)
	}
	app.Press("label:Fit view")
	app.Settle()
	// A step of the flow selects as the diagram finds it.
	pay, ok := cameraChart.Bounds("pay")
	if !ok {
		t.Fatal("the flow is not laid out")
	}
	app.PressAt("id:camera", at("flow", pay.Origin.X+pay.Size.Width/2, pay.Origin.Y+pay.Size.Height/2+28))
	if got := app.Find("id:camera-selected").Text; got != "flow/pay" {
		t.Fatalf("pressing the step selected %q", got)
	}
	app.PressAt("id:camera", driver.Point{X: 4, Y: 4})
	if got := app.Find("id:camera-selected").Text; got != "—" {
		t.Fatalf("pressing empty space selected %q", got)
	}
	before := cameraShown.ToView(core.Point{})
	app.TouchpadScroll("id:camera", driver.Point{}, driver.Point{X: 40, Y: 60})
	if after := cameraShown.ToView(core.Point{}); after.X-before.X < 39 || after.Y-before.Y < 59 {
		t.Fatalf("two fingers moved the board from %v to %v", before, after)
	}
}
