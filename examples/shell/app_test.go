package main

import (
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/nsgotest"
)

// TestMain hands the main thread to the platform under nsgo test -target.
func TestMain(m *testing.M) { nsgotest.Main(m) }

func TestShowsChosenPage(t *testing.T) {
	app := nsgotest.Start(t, Root, nsgotest.Size(960, 640))
	app.Press("id:page-sent")
	if got := app.Find("id:title").Text; got != "Sent" {
		t.Errorf("title = %q, want Sent", got)
	}
	app.Press("role:button Toggle sidebar")
	app.Advance(time.Second)
	if app.Exists("id:page-sent") {
		t.Error("the sidebar shows once hidden")
	}
	app.Press("role:button Toggle sidebar")
	app.Advance(time.Second)
	app.WaitFor("id:page-sent", time.Second)
}
