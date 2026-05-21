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
	if !strings.Contains(html, `trapFocus`) || !strings.Contains(html, `restoreFocus`) || !strings.Contains(html, `shortcutHasModifier`) || !strings.Contains(html, `gosxStudioPendingAction`) || !strings.Contains(html, `gosxStudioPendingActionLabel`) || !strings.Contains(html, `defaultPrevented`) {
		t.Fatalf("expected focus-managed command palette runtime, got: %s", html)
	}
}

func TestRenderStudioStateScriptIncludesAutosaveRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderStudioStateScript())
	if !strings.Contains(html, `data-gosx-studio-state-runtime="true"`) || !strings.Contains(html, `data-gosx-studio-state`) {
		t.Fatalf("expected embedded state runtime, got: %s", html)
	}
	for _, want := range []string{`gosxstudio:save-state`, `gosxstudio:action-result`, `gosxstudio:history-state`, `gosxstudio:history-restore`, `data-gosx-studio-client`, `data-gosx-studio-autosave`, `X-GoSX-Studio-Autosave`, `X-GoSX-Studio-Client-Action`, `data-gosx-studio-save-button`, `data-gosx-studio-last-saved`, `data-gosx-studio-dirty-count`, `data-gosx-studio-history-undo`, `data-gosx-studio-history-redo`, `actionLabel`, `requestSubmit`, `restoreFormState`} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in studio state runtime, got: %s", want, html)
		}
	}
}

func TestRenderWorkbenchScriptIncludesClientWorkbenchRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderWorkbenchScript())
	if !strings.Contains(html, `data-gosx-studio-workbench-runtime="true"`) || !strings.Contains(html, `data-studio-workbench`) {
		t.Fatalf("expected embedded workbench runtime, got: %s", html)
	}
	for _, want := range []string{
		`gosxstudio:mode-change`,
		`gosxstudio:viewport-change`,
		`gosxstudio:zoom-change`,
		`gosxstudio:rail-change`,
		`gosxstudio:workbench-rail-resize`,
		`gosxstudio:workbench-rail-resized`,
		`gosxstudio:insert-block`,
		`gosxstudio:selection-action`,
		`data-studio-rail-toggle`,
		`data-studio-resizer`,
		`--gosx-studio-left-rail-width`,
		`localStorage`,
		`setPointerCapture`,
		`data-studio-viewport`,
		`data-studio-zoom`,
		`gosxstudio:workbench-zoom`,
		`handleCommand`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in workbench runtime, got: %s", want, html)
		}
	}
}

func TestRenderFlowEditorScriptIncludesClientFlowRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderFlowEditorScript())
	if !strings.Contains(html, `data-gosx-studio-flow-editor-runtime="true"`) || !strings.Contains(html, `data-studio-flow-card`) {
		t.Fatalf("expected embedded flow editor runtime, got: %s", html)
	}
	for _, want := range []string{
		`gosxstudio:flow-select`,
		`gosxstudio:flow-dirty`,
		`gosxstudio:flow-preview`,
		`data-studio-flow-editor`,
		`data-studio-preview-frame`,
		`data-studio-open-preview`,
		`data-studio-initial-value`,
		`sessionStorage`,
		`#flow=`,
		`gosxstudio:command`,
		`gosxstudio:save-state`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in flow editor runtime, got: %s", want, html)
		}
	}
}

func TestRenderSiteCanvasScriptIncludesRuntime(t *testing.T) {
	html := gosx.RenderHTML(RenderSiteCanvasScript())
	if !strings.Contains(html, `data-gosx-studio-site-canvas-runtime="true"`) || !strings.Contains(html, `data-gosx-studio-site-canvas`) {
		t.Fatalf("expected embedded site canvas runtime, got: %s", html)
	}
	for _, want := range []string{
		`gosxstudio:canvas-select`,
		`gosxstudio:canvas-viewport`,
		`gosxstudio:canvas-cursor`,
		`gosxstudio:canvas-node-move`,
		`gosxstudio:canvas-node-moved`,
		`gosxstudio:canvas-action`,
		`data-gosx-studio-canvas-selection-detail`,
		`data-gosx-studio-canvas-node-dragging`,
		`data-gosx-studio-canvas-fit-padding`,
		`open-node`,
		`move-node`,
		`data-gosx-studio-canvas-node`,
		`gosxstudio:workbench-zoom`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in site canvas runtime, got: %s", want, html)
		}
	}
}
