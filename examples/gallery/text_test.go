package main

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
	"github.com/nsgo-dev/nsgo/nsgotest"
	"github.com/nsgo-dev/nsgo/ui"
)

//go:embed resources/fonts/Pacifico-Regular.ttf
var pacifico []byte

// TestARegisteredFontMeasures registers a font as the app runs, under nsgo test -target, and lays out a text in it with its own glyphs, not the system font's that a family it cannot find falls back to (#313).
func TestARegisteredFontMeasures(t *testing.T) {
	if nsgotest.Target() == "" {
		t.Skip("a headless test's platform registers no fonts")
	}
	const words = "Registered fonts"
	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		if err := cx.RegisterFont(pacifico); err != nil {
			t.Errorf("RegisterFont: %v", err)
		}
		return core.View{Style: core.Style{AlignItems: core.AlignStart}, Children: []core.Node{
			core.Text{TestID: "registered", Content: words, Style: core.Style{FontFamily: "Pacifico", FontSize: 20}},
			core.Text{TestID: "missing", Content: words, Style: core.Style{FontFamily: "NSGo Missing Family", FontSize: 20}},
		}}
	})
	size := func(id string) driver.Size { return app.Find("id:" + id).Frame.Size }
	// The web loads a face after RegisterFont returns.
	deadline := time.Now().Add(5 * time.Second)
	for app.Settle(); size("registered") == size("missing"); app.Settle() {
		if time.Now().After(deadline) {
			t.Fatalf("the registered font's text measures %v, as the fallback's does", size("registered"))
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestSelectingInAScrollViewKeepsItsPlace selects inside a scroll view under nsgo test -target and leaves it where it was: only a drag past its edge scrolls it, never where the last drag ended before this one's first move.
func TestSelectingInAScrollViewKeepsItsPlace(t *testing.T) {
	if nsgotest.Target() == "" {
		t.Skip("a headless drag is a finger's, which scrolls")
	}
	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		var texts []core.Node
		for i := range 40 {
			texts = append(texts, core.Text{TestID: fmt.Sprint("line-", i), Content: fmt.Sprint(i, " ", strings.Repeat("Selectable words ", 8))})
		}
		// A bar under the scroll view, as a chat's composer, puts the window's corner past its edge, and its spinner keeps the frame clock running.
		return core.View{Style: core.Style{FlexGrow: 1}, Children: []core.Node{
			core.Scroll{TestID: "scroll", Style: core.Style{FlexGrow: 1, FlexBasis: core.Pt(0)}, Children: []core.Node{
				core.View{Selectable: true, Style: core.Style{Padding: core.All(16), Gap: 12}, Children: texts},
			}},
			core.View{Style: core.Style{Height: core.Pt(40), AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{ui.Spinner{Size: 16}}},
		}}
	}, nsgotest.Size(400, 500))
	app.Scroll("id:scroll", 600)
	offset := func() core.Points { return app.Find("id:scroll").Offset.Y }
	at := offset()
	line := ""
	for i := range 40 {
		if y := app.Find(fmt.Sprint("id:line-", i)).Bounds().Origin.Y; y > 150 && y < 250 {
			line = fmt.Sprint("id:line-", i)
			break
		}
	}
	for range 2 {
		app.Drag(line, core.Point{X: 10, Y: 10}, core.Point{X: 20, Y: -20}, core.Point{X: 30, Y: -40})
		if got := offset(); got != at {
			t.Fatalf("selecting scrolled the view from %v to %v", at, got)
		}
	}
}
