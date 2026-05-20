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
}
