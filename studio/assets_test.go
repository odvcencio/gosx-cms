package studio

import (
	"strings"
	"testing"

	"github.com/odvcencio/gosx"
)

func TestRenderCommandPaletteScriptIncludesRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderCommandPaletteScript())
	if !strings.Contains(html, `data-gosx-studio-command-runtime="true"`) || !strings.Contains(html, `gosxstudio:command`) {
		t.Fatalf("expected embedded command runtime, got: %s", html)
	}
	if !strings.Contains(html, `data-studio-command-shortcut`) || !strings.Contains(html, `shortcutMatches`) {
		t.Fatalf("expected embedded command shortcuts, got: %s", html)
	}
	if !strings.Contains(html, `trapFocus`) || !strings.Contains(html, `restoreFocus`) || !strings.Contains(html, `shortcutHasModifier`) {
		t.Fatalf("expected focus-managed command palette runtime, got: %s", html)
	}
}

func TestRenderStudioStateScriptIncludesAutosaveRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderStudioStateScript())
	if !strings.Contains(html, `data-gosx-studio-state-runtime="true"`) || !strings.Contains(html, `data-gosx-studio-state`) {
		t.Fatalf("expected embedded state runtime, got: %s", html)
	}
	for _, want := range []string{`gosxstudio:save-state`, `data-gosx-studio-autosave`, `X-GoSX-Studio-Autosave`, `data-gosx-studio-save-button`, `data-gosx-studio-last-saved`, `data-gosx-studio-dirty-count`, `requestSubmit`} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in studio state runtime, got: %s", want, html)
		}
	}
}

func TestRenderSiteCanvasScriptIncludesRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderSiteCanvasScript())
	if !strings.Contains(html, `data-gosx-studio-site-canvas-runtime="true"`) || !strings.Contains(html, `data-gosx-studio-site-canvas`) {
		t.Fatalf("expected embedded site canvas runtime, got: %s", html)
	}
	for _, want := range []string{`gosxstudio:canvas-select`, `gosxstudio:canvas-viewport`, `gosxstudio:canvas-cursor`, `data-gosx-studio-canvas-node`} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in site canvas runtime, got: %s", want, html)
		}
	}
}
