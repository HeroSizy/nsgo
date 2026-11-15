// Package flow is base/flow's styled form: a directed diagram in ui's
// look, its lines in the muted text's color, the hovered and selected
// node's in the accent, and Card, a node's default face. Diagram is the
// diagram alone, for an app's own camera.View; View puts it on a camera
// of its own.
package flow

import (
	"github.com/nsgo-dev/nsgo/base/camera"
	baseflow "github.com/nsgo-dev/nsgo/base/flow"
	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/ui"
)

// The diagram's data and state are base/flow's.
type (
	Direction = baseflow.Direction
	Node      = baseflow.Node
	Edge      = baseflow.Edge
	Data      = baseflow.Data
	State     = baseflow.State
)

const (
	TopToBottom = baseflow.TopToBottom
	LeftToRight = baseflow.LeftToRight
)

// Diagram is base/flow's View in ui's look, laid out at its size; its
// fields are base's. While there are no nodes Empty shows, else
// EmptyText, an em dash when empty, in the muted text's color.
type Diagram struct {
	Data             func() Data
	Direction        func() Direction
	NodeGap, RankGap core.Points
	Selected         func() string
	OnSelect         func(id string)
	OnOpen           func(id string)
	State            *State
	Camera           *camera.State
	Origin           func() core.Point
	Empty            core.Node
	EmptyText        string
	Label            string
	LabelFunc        func() string
	Style            core.Style
	TestID           string
}

// Render implements core.Node.
func (d Diagram) Render(cx *core.Scope) core.Node {
	colors := ui.ColorsOf(cx)
	empty := d.Empty
	if empty == nil {
		text := d.EmptyText
		if text == "" {
			text = "—"
		}
		empty = core.Text{Content: text, Style: ui.TypeBody, StyleFunc: ui.Colored(colors.Muted)}
	}
	return baseflow.View{Data: d.Data, Direction: d.Direction, NodeGap: d.NodeGap, RankGap: d.RankGap,
		Selected: d.Selected, OnSelect: d.OnSelect, OnOpen: d.OnOpen, State: d.State, Camera: d.Camera, Origin: d.Origin, Empty: empty,
		Label: d.Label, LabelFunc: d.LabelFunc, Style: d.Style, TestID: d.TestID,
		LookFunc: func() baseflow.Look {
			return baseflow.Look{Line: ui.Mix(colors.Muted(), colors.Background(), .45), Emphasis: colors.Accent().Faded(.85),
				Label: colors.Muted(), LabelBackground: colors.Background(), LabelBorder: colors.Border().Faded(.7)}
		},
		NodeStyleFunc: func() core.Style {
			return core.Style{CornerRadius: CardRadius, FocusVisible: ui.FocusRing(cx, 2)}
		},
		RingStyleFunc: func() core.Style {
			return core.Style{BorderWidth: 2, BorderColor: colors.Accent().Faded(.7), CornerRadius: CardRadius + 3}
		}}
}

// overlayInset keeps the overlay from the view's corner.
const overlayInset = 8

// View is the Diagram on a camera.View of its own, which fills what the
// app sizes: it pans and zooms, presses and hovers reach the nodes, and
// Overlay sits in its top trailing corner, such as zoom controls. Camera
// holds the camera between mounts and moves it; nil keeps one of the
// view's own.
type View struct {
	Data             func() Data
	Direction        func() Direction
	NodeGap, RankGap core.Points
	Selected         func() string
	OnSelect         func(id string)
	OnOpen           func(id string)
	State            *State
	Camera           *camera.State
	Empty, Overlay   core.Node
	EmptyText        string
	Label            string
	LabelFunc        func() string
	Style            core.Style
	TestID           string
}

// Render implements core.Node.
func (v View) Render(cx *core.Scope) core.Node {
	state, cam := v.State, v.Camera
	if state == nil {
		state = &State{}
	}
	if cam == nil {
		cam = &camera.State{}
	}
	var children []core.Node
	if v.Overlay != nil {
		children = append(children, core.View{Style: core.Style{Position: core.Absolute, Top: core.Pt(overlayInset), Right: core.Pt(overlayInset)},
			Children: []core.Node{v.Overlay}})
	}
	return camera.View{State: cam, TestID: v.TestID, Style: v.Style,
		Content: Diagram{Data: v.Data, Direction: v.Direction, NodeGap: v.NodeGap, RankGap: v.RankGap, Selected: v.Selected,
			OnSelect: v.OnSelect, OnOpen: v.OnOpen, State: state, Camera: cam, Empty: v.Empty, EmptyText: v.EmptyText,
			Label: v.Label, LabelFunc: v.LabelFunc},
		OnPress: state.Press,
		OnHover: state.Hover,
		Cursor: func() core.Cursor {
			if state.Hovered() != "" {
				// Nodes are pressables, which keep the platform's arrow.
				return core.CursorDefault
			}
			return core.CursorAuto
		},
		Children: children}
}
