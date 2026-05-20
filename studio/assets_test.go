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
