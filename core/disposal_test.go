package core

import (
	goruntime "runtime"
	"testing"
	"weak"

	"github.com/nsgo-dev/nsgo/internal/driver"
)

type disposalTextEngine struct{ texts []*disposalPreparedText }

func (e *disposalTextEngine) Prepare(at *driver.AttributedText) driver.PreparedText {
	p := &disposalPreparedText{PreparedText: ApproxTextEngine{}.Prepare(at)}
	e.texts = append(e.texts, p)
	return p
}

type disposalPreparedText struct {
	driver.PreparedText
	releases int
}

func (p *disposalPreparedText) Release() {
	p.releases++
	p.PreparedText.Release()
}

func TestDisposingATreeReleasesItsPreparedTexts(t *testing.T) {
	text := &disposalTextEngine{}
	root := testRoot(RootOptions{})
	tree := NewTree(root, text)
	tree.Resize(Size{Width: 200, Height: 200}, Geometry{})
	var value *Signal[string]
	tree.Mount(func(cx *Scope) Node {
		value = NewSignal(cx, "First text")
		return View{Children: []Node{
			Text{ContentFunc: value.Get},
			Input{Value: value.Get, Multiline: true, OnContentSizeChange: func(Size) {}},
		}}
	})
	tree.Commit()
	value.Set("Replacement text")
	tree.Commit()
	if len(text.texts) == 0 {
		t.Fatal("the tree prepared no text")
	}
	root.Dispose()
	root.Dispose()
	for i, p := range text.texts {
		if p.releases != 1 {
			t.Errorf("prepared text %d released %d times, want once", i, p.releases)
		}
	}
}

// Overlays and the views holding them are freed once removed: their host lets go of them, so a tooltip or popover shown and hidden again and again does not keep each subtree it had.
func TestRemovedOverlaysAreFreed(t *testing.T) {
	var shown *Signal[bool]
	h := newHost(t, Size{Width: 400, Height: 400}, func(cx *Scope) Node {
		shown = NewSignal(cx, true)
		return View{Children: []Node{
			Show{When: shown.Get, Then: func(*Scope) Node {
				return Group{
					View{TestID: "overlaid", Children: []Node{
						Overlay{TestID: "overlay", Children: []Node{Text{Content: "Over"}}},
					}},
					View{TestID: "tipped", Tooltip: Text{Content: "Tip"}},
					Modal{TestID: "modal", Children: []Node{Text{Content: "Modal"}}},
				}
			}},
		}}
	})
	h.commit()
	var gone []weak.Pointer[node]
	for _, id := range []string{"overlaid", "overlay", "tipped", "modal"} {
		gone = append(gone, weak.Make(h.find(id)))
	}
	for _, on := range []bool{false, true, false} {
		shown.Set(on)
		h.commit()
	}
	goruntime.GC()
	for i, p := range gone {
		if p.Value() != nil {
			t.Errorf("view %d stays in memory once removed", i)
		}
	}
	if n := len(h.tree.view.hosted); n != 0 {
		t.Errorf("the root view hosts %d views once they are removed", n)
	}
	goruntime.KeepAlive(h.tree)
}
