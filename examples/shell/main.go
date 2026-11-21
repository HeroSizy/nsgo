// Command shell is a desktop window in ui.SidebarShell's inset variant:
// the content floats as a card beside a sidebar over the window's
// material, and slides over the sidebar as the toggle hides it.
package main

import (
	"github.com/nsgo-dev/nsgo"
	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/icons/lucide"
	"github.com/nsgo-dev/nsgo/ui"
)

func main() {
	nsgo.Run(Root, nsgo.WindowSize(960, 640), nsgo.WindowMaterial(core.WindowMaterialSidebar))
}

// pages are the sidebar's items, each a page of the content.
var pages = []struct {
	id, title string
	icon      func() *core.SVG
}{
	{"inbox", "Inbox", lucide.Inbox},
	{"drafts", "Drafts", lucide.FileText},
	{"sent", "Sent", lucide.Send},
	{"archive", "Archive", lucide.Archive},
}

// Root is the window: the sidebar of pages beside the chosen one.
func Root(cx *core.Scope) core.Node {
	open := core.NewSignal(cx, true)
	page := core.NewSignal(cx, pages[0].id)
	items := make([]ui.SidebarItem, len(pages))
	for i, p := range pages {
		items[i] = ui.SidebarItem{ID: p.id, Text: p.title, Icon: p.icon, TestID: "page-" + p.id}
	}
	title := func() string {
		for _, p := range pages {
			if p.id == page.Get() {
				return p.title
			}
		}
		return ""
	}
	return ui.SidebarShell{TestID: "shell", Variant: ui.SidebarShellInset,
		Open: open.Get, OnOpenChange: open.Set, ToggleLabel: "Toggle sidebar", ResizeLabel: "Resize sidebar",
		Sidebar: []core.Node{ui.Sidebar{TestID: "sidebar", Selected: page.Get, OnSelect: page.Set,
			Sections: []ui.SidebarSection{{Title: "Mail", Items: items}}}},
		Header: []core.Node{core.Text{TestID: "title", ContentFunc: title, Style: core.Style{FontWeight: core.FontWeightSemibold}}},
		Content: []core.Node{core.View{Style: core.Style{FlexGrow: 1, AlignItems: core.AlignCenter, JustifyContent: core.JustifyCenter}, Children: []core.Node{
			core.Text{TestID: "body", ContentFunc: title},
		}}}}
}
