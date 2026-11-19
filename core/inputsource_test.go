package core_test

import (
	"testing"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/nsgotest"
)

// A field's input source is selected each time it takes the focus and
// restored as it loses it, unless it keeps it.
func TestInputSelectsItsInputSource(t *testing.T) {
	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		return core.View{Children: []core.Node{
			core.Input{TestID: "code", InputSource: "abc"},
			core.Input{TestID: "name"},
			core.Input{TestID: "kept", InputSource: "xyz", KeepInputSource: true},
		}}
	})
	if app.InputSource != "" {
		t.Fatalf("selected %q before the focus", app.InputSource)
	}
	for range 2 {
		app.Focus("id:code")
		if app.InputSource != "abc" {
			t.Fatalf("focused, selected %q", app.InputSource)
		}
		app.Focus("id:name")
		if app.InputSource != "" {
			t.Fatalf("blurred, left %q", app.InputSource)
		}
	}
	app.Focus("id:kept")
	app.Focus("id:name")
	if app.InputSource != "xyz" {
		t.Fatalf("kept %q", app.InputSource)
	}
}
