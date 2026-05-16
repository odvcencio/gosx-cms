package studio

import (
	"strings"
	"testing"

	"github.com/odvcencio/gosx"
	"github.com/odvcencio/gosx-admin/blockstudio"
	"github.com/odvcencio/gosx-cms/lifecycle"
	"github.com/odvcencio/gosx-cms/media"
)

func TestNewNormalizesShell(t *testing.T) {
	shell := New(Options{
		SaveAction:   "/admin/editor?action=save",
		BlockCatalog: []blockstudio.Definition{{Key: "hero", Label: "Hero"}},
		Media:        []media.Asset{{URL: "/media/forest.jpg"}},
		Revisions:    []lifecycle.Revision{{ID: "rev_1"}},
		Left:         []Panel{{Key: "Site Map", Label: "Site map"}},
		Right:        []Panel{{Key: "Inspector", Label: "Inspector"}},
	})
	if shell.Title != "GoSX Studio" || shell.PreviewURL != "/" || shell.BlockCount != 1 || shell.MediaCount != 1 || shell.RevisionCount != 1 {
		t.Fatalf("unexpected shell: %#v", shell)
	}
	if len(shell.Actions) != 1 || shell.Actions[0].Key != "save" || !shell.Actions[0].Primary {
		t.Fatalf("expected default save action, got %#v", shell.Actions)
	}
	if shell.Left[0].Key != "site-map" || shell.Right[0].Key != "inspector" {
		t.Fatalf("expected normalized panels, got %#v %#v", shell.Left, shell.Right)
	}
}

func TestRenderShell(t *testing.T) {
	shell := New(Options{
		Title:      "Nature School Studio",
		PreviewURL: "/programs",
		Actions:    []Action{{Key: "preview", Label: "Preview", Href: "/programs"}},
		Left:       []Panel{{Key: "calendar", Label: "Calendar", Summary: "Upcoming sessions"}},
		Right:      []Panel{{Key: "inspector", Label: "Inspector", Children: []gosx.Node{gosx.El("p", nil, gosx.Text("Fields"))}}},
	})
	html := gosx.RenderHTML(Render(shell))
	for _, want := range []string{
		`class="gosx-studio"`,
		"Nature School Studio",
		`src="/programs"`,
		`data-panel-key="calendar"`,
		`data-panel-key="inspector"`,
		`data-action-key="preview"`,
		"Fields",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in shell html: %s", want, html)
		}
	}
}

func TestView(t *testing.T) {
	shell := New(Options{
		Title:         "Studio",
		PreviewURL:    "/",
		SaveAction:    "/admin/editor?action=save",
		RestoreAction: "/admin/editor?action=restoreRevision",
		BlockCatalog:  []blockstudio.Definition{{Key: "hero"}},
		Media:         []media.Asset{{URL: "/media/forest.jpg"}},
		Revisions:     []lifecycle.Revision{{ID: "rev_1"}},
		Left:          []Panel{{Key: "map", Label: "Map", Summary: "Pages"}},
		Actions:       []Action{{Key: "preview", Label: "Preview", Href: "/", Primary: false}},
	})
	view := View(shell)
	if view["title"] != "Studio" || view["saveAction"] != "/admin/editor?action=save" || view["blockCount"] != 1 || view["hasMedia"] != true || view["hasRevisions"] != true {
		t.Fatalf("unexpected view: %#v", view)
	}
	actions := view["actions"].([]map[string]any)
	if len(actions) != 2 || actions[0]["key"] != "save" || actions[0]["class"] != "button button--primary" {
		t.Fatalf("unexpected action views: %#v", actions)
	}
	left := view["leftPanels"].([]map[string]any)
	if len(left) != 1 || left[0]["key"] != "map" || left[0]["hasSummary"] != true {
		t.Fatalf("unexpected panel views: %#v", left)
	}
}
