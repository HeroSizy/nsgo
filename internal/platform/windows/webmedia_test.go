//go:build windows

package windows

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
)

// TestWebPagesCaptureWhatMediaAllows has a page of a web view that may
// capture the microphone ask for it and the camera, with Chromium's fake
// devices, and a page of one that may capture nothing ask for the
// microphone, neither prompting. It runs in a process of its own, with a
// profile of its own, whose browser process takes the fake devices' flag.
func TestWebPagesCaptureWhatMediaAllows(t *testing.T) {
	if testing.Short() {
		t.Skip("opens a window")
	}
	if _, err := findWebView2(); err != nil {
		t.Skip(err)
	}
	if os.Getenv("NSGO_TEST_WEBVIEW2_PROCESS") != t.Name() {
		cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v")
		dir := t.TempDir()
		cmd.Env = append(os.Environ(), "NSGO_TEST_WEBVIEW2_PROCESS="+t.Name(), "APPDATA="+dir, "LOCALAPPDATA="+dir,
			"WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--use-fake-device-for-media-stream")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
		if strings.Contains(string(out), "--- SKIP") {
			t.Skipf("%s", out)
		}
		return
	}
	// The page asks when the app tells it to: the view that may capture
	// nothing asks once the other has its answer, which WebView2 would
	// keep for the origin.
	const page = `<!doctype html><html><body><script>
		const ask = (what) => navigator.mediaDevices.getUserMedia(what).then((s) => s.getTracks().length, (e) => e.name);
		addEventListener('message', async () => window.nsgo.postMessage((await ask({ audio: true })) + ',' + (await ask({ video: true }))));
		window.nsgo.postMessage('ready');
	</script></body></html>`
	title := "nsgo-web-media"
	app := func(cx *core.Scope) core.Node {
		var mic, none core.Ref
		heard := map[*core.Ref]*core.Signal[string]{&mic: core.NewSignal(cx, ""), &none: core.NewSignal(cx, "")}
		cx.Effect(func() { cx.SetTitle(title + "|" + heard[&mic].Get() + "|" + heard[&none].Get()) })
		web := func(ref *core.Ref, media core.WebMedia) core.Node {
			return core.WebView{Source: core.WebHTML(page, "https://localhost/"), Media: media, Ref: ref,
				Style: core.Style{Width: driver.Pt(150), Height: driver.Pt(100)},
				OnMessage: func(data string) {
					switch {
					case data != "ready":
						heard[ref].Set(data)
						if ref == &mic {
							none.PostMessage("ask")
						}
					case ref == &mic:
						mic.PostMessage("ask")
					}
				}}
		}
		return core.View{Style: core.Style{Width: driver.Pt(320), Height: driver.Pt(240), Direction: core.Row},
			Children: []core.Node{web(&mic, core.WebMicrophone), web(&none, 0)}}
	}
	done := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		defer close(done)
		Run(app, Options{Title: title + "||", Width: 320, Height: 240})
	}()
	hwnd := waitForWindow(t, title+"||", done)
	if hwnd == 0 {
		return
	}
	t.Cleanup(func() {
		postMessageW.Call(hwnd, wmClose, 0, 0)
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("the window did not close")
		}
	})
	// The first start of the browser process can take a while.
	want := title + "|1,NotAllowedError|NotAllowedError,NotAllowedError"
	deadline := time.Now().Add(30 * time.Second)
	for got := windowText(hwnd); got != want; got = windowText(hwnd) {
		if time.Now().After(deadline) {
			t.Fatalf("window title = %q, want %q", got, want)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
