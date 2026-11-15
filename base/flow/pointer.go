package flow

import (
	"math"

	"github.com/nsgo-dev/nsgo/core"
)

// hit returns the node at a point of the diagram, or -1.
func (v *view) hit(p core.Point) int {
	if !v.ready.GetUntracked() {
		return -1
	}
	gx, gy := float64(p.X), float64(p.Y)
	for i := len(v.nodes) - 1; i >= 0; i-- {
		if i >= len(v.result.x) || !v.nodes[i].measured {
			continue
		}
		w, h := v.nodes[i].size()
		if x, y := v.result.x[i], v.result.y[i]; gx >= x && gx <= x+w && gy >= y && gy <= y+h {
			return i
		}
	}
	return -1
}

// press selects and opens the node at a point of the diagram, or selects
// "" for empty space, and says a node took it.
func (v *view) press(p core.Point, e core.PressEvent) bool {
	if i := v.hit(p); i >= 0 {
		if f := v.pressHandler(v.nodes[i].id); f != nil {
			f(e)
		}
		return true
	}
	v.lastPress = ""
	if v.spec.OnSelect != nil && v.selected() != "" {
		v.spec.OnSelect("")
	}
	return false
}

// reveal has the camera show the node with id: centered, or only when it
// does not show whole.
func (v *view) reveal(id string, center bool) {
	cam := v.spec.Camera
	b, ok := v.s.Bounds(id)
	if cam == nil || !ok {
		return
	}
	if v.spec.Origin != nil {
		o := core.Untracked(v.cx, v.spec.Origin)
		b.Origin.X += o.X
		b.Origin.Y += o.Y
	}
	if center {
		cam.Center(b)
		return
	}
	cam.Reveal(b)
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(v, hi)) }
