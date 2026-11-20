//go:build windows

package windows

import (
	"runtime"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
)

// TestClosingAWindowReleasesItsViews opens a window with a shader view and closes it: its views let go of what they draw with, the shader's texture among it, and the window of its swap chain.
func TestClosingAWindowReleasesItsViews(t *testing.T) {
	if testing.Short() {
		t.Skip("opens windows")
	}
	if _, err := devices(); err != nil {
		t.Skip(err)
	}
	const title = "nsgo-closed-views"
	prog := core.MustCompileShader(`package main

func Fragment(dst vec4, pos vec2) vec4 {
	return vec4(0, 0, 1, 1)
}
`)
	type opening struct {
		cx  *core.Scope
		win *core.Window
	}
	opened := make(chan opening, 1)
	app := core.AppRoot(func(cx *core.Scope) {
		opened <- opening{cx, cx.OpenWindow(core.WindowOptions{Title: title + "-opened", Width: 200, Height: 100}, func(*core.Scope) core.Node {
			return core.Shader{Program: prog, Style: core.Style{Width: driver.Pt(200), Height: driver.Pt(100)}}
		})}
	})
	done := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		defer close(done)
		Run(app, Options{Title: title, Width: 200, Height: 100, Hidden: true})
	}()
	root := waitForWindow(t, title, done)
	t.Cleanup(func() {
		postMessageW.Call(root, wmClose, 0, 0)
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("the app did not close")
		}
	})
	hwnd := waitForWindow(t, title+"-opened", done)
	o := <-opened
	var w *window
	var shaders []*view
	for deadline := time.Now().Add(5 * time.Second); len(shaders) == 0 || w.swap == nil; time.Sleep(20 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("the shader was not drawn")
		}
		onMain(t, hwnd, func(ow *window) {
			w, shaders = ow, nil
			for _, v := range ow.views {
				if v != nil && v.shaderView != nil && v.shaderView.bitmap != nil {
					shaders = append(shaders, v)
				}
			}
		})
	}
	o.cx.Post(o.win.Close)
	for deadline := time.Now().Add(5 * time.Second); ; time.Sleep(20 * time.Millisecond) {
		var closed bool
		onMain(t, root, func(*window) { closed = windows[hwnd] == nil })
		if closed {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the window did not close")
		}
	}
	onMain(t, root, func(*window) {
		for _, v := range shaders {
			if !v.deleted || v.shaderView != nil {
				t.Error("a closed window's shader view keeps its texture")
			}
		}
		if w.swap != nil {
			t.Error("a closed window keeps its swap chain")
		}
	})
}
