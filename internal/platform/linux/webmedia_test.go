//go:build linux && !android

package linux

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/internal/driver"
)

// TestWebPagesCaptureWhatMediaAllows has a page of a web view that may
// capture the microphone ask for it and the camera, with WebKit's mock
// devices, and a page of one that may capture nothing find no devices.
// It skips where WebKitGTK 6.0 is missing, in a process of its own as
// WebKit takes the first thread that uses it.
func TestWebPagesCaptureWhatMediaAllows(t *testing.T) {
	if testing.Short() {
		t.Skip("opens a window")
	}
	if os.Getenv("NSGO_TEST_WEBKIT_PROCESS") != t.Name() {
		cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v")
		cmd.Env = append(os.Environ(), "NSGO_TEST_WEBKIT_PROCESS="+t.Name())
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
		if strings.Contains(string(out), "--- SKIP") {
			t.Skipf("%s", out)
		}
		return
	}
	if err := openWebKit(); err != nil {
		t.Skip(err)
	}
	// The page asks once the test has turned the mock devices on.
	const page = `<!doctype html><html><body><script>
		const ask = (what) => navigator.mediaDevices.getUserMedia(what).then((s) => s.getTracks().length, (e) => e.name);
		window.start = async () => {
			if (!navigator.mediaDevices) return window.nsgo.postMessage('none');
			window.nsgo.postMessage((await ask({ audio: true })) + ',' + (await ask({ video: true })));
		};
		window.nsgo.postMessage('ready');
	</script></body></html>`
	title := "nsgo-native-" + t.Name()
	app := func(cx *core.Scope) core.Node {
		heard := map[string]*core.Signal[string]{"mic": core.NewSignal(cx, ""), "none": core.NewSignal(cx, "")}
		cx.Effect(func() { cx.SetTitle(title + "|" + heard["mic"].Get() + "|" + heard["none"].Get()) })
		web := func(id string, media core.WebMedia) core.Node {
			return core.WebView{Source: core.WebHTML(page, "https://localhost/"), Media: media, TestID: id,
				Style:     core.Style{Width: driver.Pt(150), Height: driver.Pt(100)},
				OnMessage: func(data string) { heard[id].Set(data) }}
		}
		return core.View{Style: core.Style{Width: driver.Pt(360), Height: driver.Pt(240), Direction: core.Row},
			Children: []core.Node{web("mic", core.WebMicrophone), web("none", 0)}}
	}
	w := startTestWindow(t, app, Options{Title: title, Width: 360, Height: 240})
	waitForTitle(t, w, title+"|ready|ready")
	w.onLoop(t, func() {
		for _, v := range w.views {
			if v != nil && v.web != nil && v.web.widget != 0 {
				gObjectSetBool(webkitWebViewGetSettings(v.web.widget), cstr("enable-mock-capture-devices"), 1, 0)
				evaluate(v.web.widget, "start()")
			}
		}
	})
	waitForTitle(t, w, title+"|1,NotAllowedError|none")
}
