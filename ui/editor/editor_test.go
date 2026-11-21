package editor

import (
	"reflect"
	"testing"

	"github.com/nsgo-dev/nsgo/base/editor"
	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/nsgotest"
	"github.com/nsgo-dev/nsgo/ui"
)

// TestEditorTakesUIsColorsAndFollowsTheTheme mounts code in ui's colors and
// switches to the dark appearance under it.
func TestEditorTakesUIsColorsAndFollowsTheTheme(t *testing.T) {
	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		return Editor{TestID: "editor", Label: "Code", DefaultValue: "plain words\n", Style: core.Style{Width: core.Pt(400), Height: core.Pt(300)}}
	}, nsgotest.Size(600, 400))
	app.Settle()
	text := app.Find("id:editor >> text:plain words")
	light := text.TextStyle.Color
	if light == 0 || text.TextStyle.Font.Size != ui.TypeCode.FontSize {
		t.Fatalf("color %#x, size %v; want ui's text color at %v\n%s", light, text.TextStyle.Font.Size, ui.TypeCode.FontSize, app.Dump())
	}
	app.SetEnvironment(core.Environment{Dark: true})
	app.Settle()
	if dark := app.Find("id:editor >> text:plain words").TextStyle.Color; dark == 0 || dark == light {
		t.Fatalf("dark color %#x, light %#x; want the dark palette's", dark, light)
	}
}

// TestEditorMergesTheAppsStylesOverUIs gives the content the app's color,
// keeping ui's size.
func TestEditorMergesTheAppsStylesOverUIs(t *testing.T) {
	red := core.RGB(0xff0000)
	app := nsgotest.Start(t, func(cx *core.Scope) core.Node {
		return Editor{TestID: "editor", DefaultValue: "plain words\n", Styles: editor.Styles{Content: core.Style{Color: red}},
			Style: core.Style{Width: core.Pt(400), Height: core.Pt(300)}}
	}, nsgotest.Size(600, 400))
	app.Settle()
	if c := app.Find("id:editor >> text:plain words").TextStyle.Color; c != red {
		t.Fatalf("color %#x, want the app's %#x\n%s", c, red, app.Dump())
	}
}

// TestMergeCoversEveryPart sets each part in the app's styles and expects
// it over ui's.
func TestMergeCoversEveryPart(t *testing.T) {
	var app editor.Styles
	v := reflect.ValueOf(&app).Elem()
	for i := range v.NumField() {
		v.Field(i).Set(reflect.ValueOf(core.Style{FontSize: core.Points(i + 1)}))
	}
	got := reflect.ValueOf(merge(editor.Styles{}, app))
	for i := range got.NumField() {
		if s := got.Field(i).Interface().(core.Style); s.FontSize != core.Points(i+1) {
			t.Errorf("%s not merged", got.Type().Field(i).Name)
		}
	}
}
