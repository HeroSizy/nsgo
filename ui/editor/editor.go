// Package editor is base/editor's Editor in ui's colors. It is a package of
// its own, as base/editor is, so that only apps that edit code link chroma's
// lexers. Its types are base/editor's: this package adds the look alone.
package editor

import (
	"github.com/nsgo-dev/nsgo/base/editor"
	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/ui"
)

// Editor is editor.Editor in ui's look: ui.TypeCode's size and line height
// (FontSize and LineHeight replace them), faint line numbers with the
// active line's in the text's color, the accent for the caret and the
// selection, amber search matches, wavy diagnostics in the danger, warning,
// info and muted colors, and the completion list and its documentation on
// the menu's surface. The app's Styles or StylesFunc merge over ui's, part
// by part.
type Editor editor.Editor

// Render implements core.Node.
func (e Editor) Render(cx *core.Scope) core.Node {
	out := editor.Editor(e)
	if out.FontSize == 0 {
		out.FontSize = ui.TypeCode.FontSize
	}
	if out.LineHeight == 0 {
		out.LineHeight = out.FontSize * ui.TypeCode.LineHeight / ui.TypeCode.FontSize
	}
	out.Styles = editor.Styles{}
	out.StylesFunc = ui.Themed(cx, func() editor.Styles {
		app := e.Styles
		if e.StylesFunc != nil {
			app = e.StylesFunc()
		}
		return merge(Styles(cx), app)
	})
	return out
}

// Styles are Editor's styles where cx is, for an editor.Editor the app
// styles further. It reads the palette, so read it where ui.Themed does.
func Styles(cx *core.Scope) editor.Styles {
	c := ui.ColorsOf(cx).Now()
	menu := ui.FinOf(cx).Menu()
	wavy := func(color core.Color) core.Style {
		return core.Style{TextDecoration: core.TextDecorationUnderline, TextDecorationStyle: core.DecorationWavy, TextDecorationColor: color}
	}
	return editor.Styles{
		Gutter:                 core.Style{Color: c.Faint},
		ActiveLineNumber:       core.Style{Color: c.Text},
		FoldMarker:             core.Style{Color: c.Faint},
		Content:                core.Style{Color: c.Text},
		ActiveLine:             core.Style{Background: ui.Alpha(c.Text, 0.04)},
		Selection:              core.Style{Background: ui.Alpha(c.Accent, 0.22)},
		Caret:                  core.Style{Background: c.Accent, Width: core.Pt(2)},
		SearchMatch:            core.Style{Background: core.RGBA(0xfbbf24, 0.25)},
		CurrentSearchMatch:     core.Style{Background: core.RGBA(0xfbbf24, 0.55)},
		MatchingBracket:        core.Style{BorderWidth: 1, BorderColor: ui.Alpha(c.Muted, 0.6)},
		IndentGuide:            core.Style{Background: c.Border},
		FoldPlaceholder:        core.Style{Color: c.Faint},
		Placeholder:            core.Style{Color: c.Faint},
		DiagnosticError:        wavy(c.Danger),
		DiagnosticWarning:      wavy(c.Warn),
		DiagnosticInfo:         wavy(c.Info),
		DiagnosticHint:         wavy(c.Muted),
		UnnecessaryText:        core.Style{Color: c.Faint},
		DeprecatedText:         core.Style{TextDecoration: core.TextDecorationLineThrough},
		SelectionHandle:        core.Style{CornerRadius: 5, Background: c.Accent},
		CompletionMenu:         menu.Over(core.Style{CornerRadius: 8, BorderWidth: 1, BorderColor: c.Border, Padding: core.All(4)}),
		CompletionItem:         core.Style{Padding: core.Symmetric(0, 8), Gap: 6, CornerRadius: 5, Color: c.Text},
		SelectedCompletionItem: core.Style{Background: ui.Alpha(c.Accent, 0.16)},
		CompletionMatch:        core.Style{FontWeight: ui.Semibold, Color: c.Accent},
		CompletionDetail:       core.Style{Color: c.Muted, FontSize: 12},
		CompletionKind:         core.Style{Color: c.Faint, FontSize: 11, Width: core.Pt(30)},
		CompletionDocumentation: menu.Over(core.Style{CornerRadius: 8, BorderWidth: 1, BorderColor: c.Border,
			Padding: core.Symmetric(6, 8), Color: c.Text, FontSize: 12}),
		DefinitionLink:   core.Style{Background: c.Accent},
		InlineCompletion: core.Style{Color: c.Faint, Background: c.Surface},
		Whitespace:       core.Style{Color: ui.Alpha(c.Faint, 0.6)},
		Ruler:            core.Style{Background: c.Line},
	}
}

// merge is s with o over it, part by part.
func merge(s, o editor.Styles) editor.Styles {
	for _, p := range []struct{ s, o *core.Style }{
		{&s.Gutter, &o.Gutter}, {&s.LineNumber, &o.LineNumber}, {&s.ActiveLineNumber, &o.ActiveLineNumber}, {&s.FoldMarker, &o.FoldMarker},
		{&s.Content, &o.Content}, {&s.Line, &o.Line}, {&s.ActiveLine, &o.ActiveLine},
		{&s.Selection, &o.Selection}, {&s.Caret, &o.Caret}, {&s.Composition, &o.Composition},
		{&s.SearchMatch, &o.SearchMatch}, {&s.CurrentSearchMatch, &o.CurrentSearchMatch},
		{&s.MatchingBracket, &o.MatchingBracket}, {&s.IndentGuide, &o.IndentGuide}, {&s.FoldPlaceholder, &o.FoldPlaceholder},
		{&s.Placeholder, &o.Placeholder},
		{&s.DiagnosticError, &o.DiagnosticError}, {&s.DiagnosticWarning, &o.DiagnosticWarning},
		{&s.DiagnosticInfo, &o.DiagnosticInfo}, {&s.DiagnosticHint, &o.DiagnosticHint},
		{&s.UnnecessaryText, &o.UnnecessaryText}, {&s.DeprecatedText, &o.DeprecatedText},
		{&s.SelectionHandle, &o.SelectionHandle},
		{&s.CompletionMenu, &o.CompletionMenu}, {&s.CompletionItem, &o.CompletionItem}, {&s.SelectedCompletionItem, &o.SelectedCompletionItem},
		{&s.CompletionMatch, &o.CompletionMatch}, {&s.CompletionDetail, &o.CompletionDetail}, {&s.CompletionKind, &o.CompletionKind},
		{&s.CompletionDocumentation, &o.CompletionDocumentation},
		{&s.DefinitionLink, &o.DefinitionLink}, {&s.InlineCompletion, &o.InlineCompletion},
		{&s.Whitespace, &o.Whitespace}, {&s.Ruler, &o.Ruler}, {&s.OverviewRuler, &o.OverviewRuler},
	} {
		*p.s = p.s.Merge(*p.o)
	}
	return s
}
