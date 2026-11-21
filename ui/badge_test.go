package ui

import (
	"testing"

	"github.com/nsgo-dev/nsgo/core"
)

func TestBadgeVariants(t *testing.T) {
	app, p := start(t, core.Environment{}, func(*core.Scope) core.Node {
		return core.View{Children: []core.Node{
			Badge{Text: "Neutral"}, Badge{Text: "Accent", Variant: BadgeAccent}, Badge{Text: "Success", Variant: BadgeSuccess},
			Badge{Text: "Warn", Variant: BadgeWarn}, Badge{Text: "Danger", Variant: BadgeDanger},
		}}
	})
	if s := app.Find("text:Neutral").Parent.Style; s.Background != p.Surface || s.BorderColor != p.Border {
		t.Fatalf("neutral badge = %x on %x, want the surface's", s.Background, s.BorderColor)
	}
	for _, c := range []struct {
		text string
		tint core.Color
	}{{"Accent", p.Accent}, {"Success", p.Up}, {"Warn", p.Warn}, {"Danger", p.Danger}} {
		if s := app.Find("text:" + c.text).Parent.Style; s.Background != Mix(p.Background, c.tint, 0.12) || s.BorderColor != Mix(p.Background, c.tint, 0.3) {
			t.Fatalf("%s badge = %x on %x, not tinted with %x", c.text, s.Background, s.BorderColor, c.tint)
		}
	}
}
