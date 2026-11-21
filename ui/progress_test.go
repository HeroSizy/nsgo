package ui

import (
	"math"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/base"
	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
)

// TestProgress fills a determinate bar to its value and sweeps an
// indeterminate one.
func TestProgress(t *testing.T) {
	each(t, func(t *testing.T, e core.Environment) {
		app, _ := start(t, e, func(cx *core.Scope) core.Node {
			return core.View{Style: core.Style{Gap: 16}, Children: []core.Node{
				Progress{Label: "Upload", TestID: "upload", Value: func() float32 { return 40 }},
				Progress{Label: "Sync", TestID: "sync", Indeterminate: func() bool { return true }},
			}}
		})
		bar := app.Find("upload")
		track, fill := bar.Children[0], bar.Children[0].Children[0]
		if r := fill.Frame.Size.Width / track.Frame.Size.Width; math.Abs(float64(r)-0.4) > 0.01 {
			t.Fatalf("fill = %v of the track, want 0.4", r)
		}
		if bar.Props.Value.Now != 40 {
			t.Fatalf("value = %v, want 40", bar.Props.Value.Now)
		}
		sweep := func() core.Points { return app.Find("sync").Children[0].Children[0].Style.Transform.TranslateX }
		before := sweep()
		app.FastForward(time.Second)
		if sweep() == before {
			t.Fatal("the indeterminate bar does not sweep")
		}
	})
}

// TestProgressDirections fills a bar from the end Axis and Inverted say.
func TestProgressDirections(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		forty := func() float32 { return 40 }
		return core.View{Style: core.Style{Direction: core.Row, Gap: 16, AlignItems: core.AlignStart}, Children: []core.Node{
			Progress{TestID: "up", Value: forty, Axis: base.Vertical, Style: core.Style{Height: core.Pt(100)}},
			Progress{TestID: "down", Value: forty, Axis: base.Vertical, Inverted: true, Style: core.Style{Height: core.Pt(100)}},
			Progress{TestID: "back", Value: forty, Inverted: true, Style: core.Style{Width: core.Pt(100)}},
		}}
	})
	fill := func(id string) driver.Bounds { return app.Find(id).Children[0].Children[0].Frame }
	if f := fill("up"); f.Size.Height != 40 || f.Origin.Y != 60 || f.Size.Width != 6 {
		t.Errorf("bottom-up fill = %v, want 6 × 40 at the bottom", f)
	}
	if f := fill("down"); f.Size.Height != 40 || f.Origin.Y != 0 {
		t.Errorf("top-down fill = %v, want 40 high at the top", f)
	}
	if f := fill("back"); f.Size.Width != 40 || f.Origin.X != 60 {
		t.Errorf("inverted fill = %v, want 40 wide at the end", f)
	}
}
