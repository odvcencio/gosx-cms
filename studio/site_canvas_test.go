package studio

import (
	"strings"
	"testing"

	"github.com/odvcencio/gosx"
)

func TestRenderSiteCanvasMapsWebsiteObjects(t *testing.T) {
	html := gosx.RenderHTML(RenderSiteCanvas(SiteCanvasOptions{
		Class:         "studio-site-map",
		ControlsClass: "studio-site-map__controls",
		NodeClass:     "studio-site-map__node",
		Kicker:        "Spatial authoring",
		Title:         "Website map",
		Summary:       "Pages, flows, content, and style",
		Nodes: []SiteCanvasNode{
			{Key: "home", Kind: "page", Label: "Home", Summary: "Landing page", Href: "/", X: 120, Y: 120, Selected: true, Metrics: []Metric{NewMetric("sections", "sections", 6)}},
			{Key: "checkout", Kind: "flow", Label: "Checkout", Summary: "Purchase path", Href: "/shop", X: 480, Y: 180},
			{Key: "theme", Kind: "style", Label: "Theme", Summary: "Palette and type", X: 120, Y: 360},
		},
		Edges: []SiteCanvasEdge{
			{From: "home", To: "checkout", Kind: "journey", Label: "Buy path"},
			{From: "theme", To: "home", Kind: "style"},
		},
	}))

	for _, want := range []string{
		`data-gosx-studio-site-canvas="true"`,
		`data-gosx-studio-canvas-viewport="true"`,
		`data-gosx-studio-canvas-surface="true"`,
		`data-gosx-studio-canvas-node="home"`,
		`data-gosx-studio-canvas-node-kind="style"`,
		`data-gosx-studio-canvas-edge-from="home"`,
		`data-gosx-studio-canvas-edge-to="checkout"`,
		`class="studio-site-map__controls"`,
		`studio-site-map__node--page`,
		`studio-site-map__node-kind`,
		`studio-site-map__node-metrics`,
		`aria-pressed="true"`,
		`Website map`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in site canvas markup: %s", want, html)
		}
	}
}
