package studio

import (
	"strings"
	"testing"
	"time"

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

func TestNewDefaults(t *testing.T) {
	shell := New(Options{})
	if shell.Title != "GoSX Studio" || shell.PreviewURL != "/" {
		t.Fatalf("unexpected defaults: %#v", shell)
	}
	if shell.SaveAction != "" || shell.RestoreAction != "" {
		t.Fatalf("expected empty actions URLs by default: %#v", shell)
	}
	if shell.BlockCount != 0 || shell.MediaCount != 0 || shell.RevisionCount != 0 || shell.HasMedia || shell.HasRevisions {
		t.Fatalf("unexpected empty counts: %#v", shell)
	}
}

func TestActionNormalization(t *testing.T) {
	shell := New(Options{
		SaveAction: "/save",
		Actions: []Action{
			{Key: "Preview Page", Label: " Preview ", Href: " /preview "},
			{Key: "save", Label: "Custom save", Href: "/custom-save"},
			{Key: "missing-label", Href: "/skip"},
		},
		Navigation: []Section{{
			Key:     "Main Nav",
			Label:   " Main ",
			Summary: " Primary links ",
			Actions: []Action{
				{Key: "Media Library", Label: " Media ", Href: " /media "},
				{Key: "invalid"},
			},
		}},
	})
	if len(shell.Actions) != 2 {
		t.Fatalf("expected custom save to prevent default save insertion: %#v", shell.Actions)
	}
	if shell.Actions[0].Key != "preview-page" || shell.Actions[0].Label != "Preview" || shell.Actions[0].Href != "/preview" {
		t.Fatalf("unexpected normalized preview action: %#v", shell.Actions[0])
	}
	if shell.Actions[1].Key != "save" || shell.Actions[1].Label != "Custom save" || shell.Actions[1].Href != "/custom-save" {
		t.Fatalf("unexpected normalized save action: %#v", shell.Actions[1])
	}
	if len(shell.Navigation) != 1 || shell.Navigation[0].Key != "main-nav" || shell.Navigation[0].Summary != "Primary links" {
		t.Fatalf("unexpected navigation section: %#v", shell.Navigation)
	}
	if len(shell.Navigation[0].Actions) != 1 || shell.Navigation[0].Actions[0].Key != "media-library" {
		t.Fatalf("unexpected navigation actions: %#v", shell.Navigation[0].Actions)
	}
}

func TestPanelAndCatalogSummaries(t *testing.T) {
	archivedAt := time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC)
	created := time.Date(2026, 5, 2, 10, 30, 0, 0, time.UTC)
	shell := New(Options{
		BlockCatalog: []blockstudio.Definition{{
			Key:        "hero",
			Label:      "Hero",
			Summary:    "Lead story",
			Kind:       "section",
			Preview:    "/blocks/hero",
			DefaultOn:  true,
			Locked:     true,
			Repeatable: true,
			Icon:       "image",
			Fields:     []blockstudio.FieldDefinition{{Name: "title"}},
		}},
		Media: []media.Asset{{
			ID:          "asset_1",
			URL:         "/media/forest.jpg",
			Alt:         "Forest",
			Filename:    "forest.jpg",
			ContentType: "image/jpeg",
			Size:        2048,
			ArchivedAt:  &archivedAt,
		}},
		Revisions: []lifecycle.Revision{{
			ID:            "rev_1",
			ResourceKind:  "page",
			ResourceID:    "home",
			ResourceTitle: "Home",
			Action:        "publish",
			Summary:       "Published home",
			Created:       created,
		}},
		Left:  []Panel{NewPanel("Content Tree", "Content tree", "3 pages")},
		Right: []Panel{{Key: " Inspector ", Label: " Inspector ", Summary: " Fields "}},
	})
	if len(shell.Blocks) != 1 || shell.Blocks[0].FieldCount != 1 || !shell.Blocks[0].DefaultOn || !shell.Blocks[0].Locked || !shell.Blocks[0].Repeatable {
		t.Fatalf("unexpected block summary: %#v", shell.Blocks)
	}
	if len(shell.MediaLibrary) != 1 || !shell.MediaLibrary[0].Archived || shell.MediaLibrary[0].Filename != "forest.jpg" {
		t.Fatalf("unexpected media summary: %#v", shell.MediaLibrary)
	}
	if len(shell.RevisionLog) != 1 || shell.RevisionLog[0].Created != "2026-05-02T10:30:00Z" {
		t.Fatalf("unexpected revision summary: %#v", shell.RevisionLog)
	}
	if shell.Left[0].Key != "content-tree" || shell.Left[0].Summary != "3 pages" || shell.Right[0].Key != "inspector" || shell.Right[0].Summary != "Fields" {
		t.Fatalf("unexpected panel summaries: %#v %#v", shell.Left, shell.Right)
	}
}

func TestMetricsAndExtras(t *testing.T) {
	extras := map[string]any{
		"workflow": "draft",
		"title":    "should not replace core title",
	}
	shell := New(Options{
		Title:   "Studio",
		Metrics: []Metric{NewMetric("Open Drafts", "Open drafts", 4), {Key: "skip"}},
		Extras:  extras,
	})
	extras["workflow"] = "published"
	if len(shell.Metrics) != 1 || shell.Metrics[0].Key != "open-drafts" || shell.Metrics[0].Value != 4 {
		t.Fatalf("unexpected metrics: %#v", shell.Metrics)
	}
	if shell.Extras["workflow"] != "draft" {
		t.Fatalf("expected extras to be cloned: %#v", shell.Extras)
	}
	view := View(shell)
	if view["title"] != "Studio" || view["workflow"] != "draft" {
		t.Fatalf("unexpected extras passthrough: %#v", view)
	}
	metrics := view["metrics"].([]map[string]any)
	if len(metrics) != 1 || metrics[0]["key"] != "open-drafts" || metrics[0]["value"] != 4 {
		t.Fatalf("unexpected metric views: %#v", metrics)
	}
	viewExtras := view["extras"].(map[string]any)
	viewExtras["workflow"] = "changed"
	if shell.Extras["workflow"] != "draft" {
		t.Fatalf("expected view extras to be cloned: %#v", shell.Extras)
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
	for _, key := range []string{"previewURL", "restoreAction", "mediaCount", "revisionCount", "actions", "leftPanels", "rightPanels"} {
		if _, ok := view[key]; !ok {
			t.Fatalf("expected compatibility key %q in view: %#v", key, view)
		}
	}
	actions := view["actions"].([]map[string]any)
	if len(actions) != 2 || actions[0]["key"] != "save" || actions[0]["class"] != "button button--primary" {
		t.Fatalf("unexpected action views: %#v", actions)
	}
	left := view["leftPanels"].([]map[string]any)
	if len(left) != 1 || left[0]["key"] != "map" || left[0]["hasSummary"] != true {
		t.Fatalf("unexpected panel views: %#v", left)
	}
	blocks := view["blockCatalog"].([]map[string]any)
	media := view["media"].([]map[string]any)
	revisions := view["revisions"].([]map[string]any)
	if len(blocks) != 1 || blocks[0]["key"] != "hero" || len(media) != 1 || media[0]["url"] != "/media/forest.jpg" || len(revisions) != 1 || revisions[0]["id"] != "rev_1" {
		t.Fatalf("unexpected catalog views: %#v %#v %#v", blocks, media, revisions)
	}
}
