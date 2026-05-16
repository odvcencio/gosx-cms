package studio

import (
	"strings"

	"github.com/odvcencio/gosx"
	"github.com/odvcencio/gosx-admin/blockstudio"
	"github.com/odvcencio/gosx-cms/lifecycle"
	"github.com/odvcencio/gosx-cms/media"
)

type Options struct {
	Title         string
	PreviewURL    string
	SaveAction    string
	RestoreAction string
	BlockCatalog  []blockstudio.Definition
	Media         []media.Asset
	Revisions     []lifecycle.Revision
	Left          []Panel
	Right         []Panel
	Actions       []Action
}

type Action struct {
	Key     string
	Label   string
	Href    string
	Primary bool
}

type Panel struct {
	Key      string
	Label    string
	Summary  string
	Children []gosx.Node
}

type Shell struct {
	Title         string
	PreviewURL    string
	SaveAction    string
	RestoreAction string
	BlockCatalog  []blockstudio.Definition
	BlockCount    int
	Media         []media.Asset
	MediaCount    int
	HasMedia      bool
	Revisions     []lifecycle.Revision
	RevisionCount int
	HasRevisions  bool
	Left          []Panel
	Right         []Panel
	Actions       []Action
}

func New(options Options) Shell {
	title := strings.TrimSpace(options.Title)
	if title == "" {
		title = "GoSX Studio"
	}
	previewURL := strings.TrimSpace(options.PreviewURL)
	if previewURL == "" {
		previewURL = "/"
	}
	actions := normalizeActions(options.Actions, options.SaveAction)
	mediaAssets := media.CloneAssets(options.Media)
	revisions := lifecycle.CloneRevisions(options.Revisions)
	return Shell{
		Title:         title,
		PreviewURL:    previewURL,
		SaveAction:    strings.TrimSpace(options.SaveAction),
		RestoreAction: strings.TrimSpace(options.RestoreAction),
		BlockCatalog:  append([]blockstudio.Definition(nil), options.BlockCatalog...),
		BlockCount:    len(options.BlockCatalog),
		Media:         mediaAssets,
		MediaCount:    len(mediaAssets),
		HasMedia:      len(mediaAssets) > 0,
		Revisions:     revisions,
		RevisionCount: len(revisions),
		HasRevisions:  len(revisions) > 0,
		Left:          normalizePanels(options.Left),
		Right:         normalizePanels(options.Right),
		Actions:       actions,
	}
}

func Render(shell Shell) gosx.Node {
	return gosx.El("div", gosx.Attrs(gosx.Attr("class", "gosx-studio")),
		gosx.El("header", gosx.Attrs(gosx.Attr("class", "gosx-studio__toolbar")),
			gosx.El("h1", nil, gosx.Text(shell.Title)),
			renderActions(shell.Actions),
		),
		gosx.El("div", gosx.Attrs(gosx.Attr("class", "gosx-studio__layout")),
			renderPanelColumn("gosx-studio__left", shell.Left),
			gosx.El("main", gosx.Attrs(gosx.Attr("class", "gosx-studio__canvas")),
				gosx.El("iframe", gosx.Attrs(
					gosx.Attr("src", shell.PreviewURL),
					gosx.Attr("title", shell.Title+" preview"),
					gosx.Attr("data-studio-preview", "true"),
				)),
			),
			renderPanelColumn("gosx-studio__right", shell.Right),
		),
	)
}

func renderActions(actions []Action) gosx.Node {
	nodes := make([]gosx.Node, 0, len(actions))
	for _, action := range actions {
		className := "button button--secondary"
		if action.Primary {
			className = "button button--primary"
		}
		nodes = append(nodes, gosx.El("a", gosx.Attrs(
			gosx.Attr("class", className),
			gosx.Attr("href", action.Href),
			gosx.Attr("data-action-key", action.Key),
		), gosx.Text(action.Label)))
	}
	return gosx.El("nav", gosx.Attrs(gosx.Attr("class", "gosx-studio__actions")), gosx.Fragment(nodes...))
}

func renderPanelColumn(className string, panels []Panel) gosx.Node {
	nodes := make([]gosx.Node, 0, len(panels))
	for _, panel := range panels {
		children := []gosx.Node{gosx.El("h2", nil, gosx.Text(panel.Label))}
		if panel.Summary != "" {
			children = append(children, gosx.El("p", nil, gosx.Text(panel.Summary)))
		}
		children = append(children, panel.Children...)
		nodes = append(nodes, gosx.El("section", gosx.Attrs(
			gosx.Attr("class", "gosx-studio__panel"),
			gosx.Attr("data-panel-key", panel.Key),
		), gosx.Fragment(children...)))
	}
	return gosx.El("aside", gosx.Attrs(gosx.Attr("class", className)), gosx.Fragment(nodes...))
}

func normalizeActions(actions []Action, saveAction string) []Action {
	out := make([]Action, 0, len(actions)+1)
	for _, action := range actions {
		action.Key = normalizeKey(action.Key)
		action.Label = strings.TrimSpace(action.Label)
		action.Href = strings.TrimSpace(action.Href)
		if action.Key == "" || action.Label == "" {
			continue
		}
		out = append(out, action)
	}
	if strings.TrimSpace(saveAction) != "" && !hasAction(out, "save") {
		out = append([]Action{{Key: "save", Label: "Save", Href: strings.TrimSpace(saveAction), Primary: true}}, out...)
	}
	return out
}

func normalizePanels(panels []Panel) []Panel {
	out := make([]Panel, 0, len(panels))
	for _, panel := range panels {
		panel.Key = normalizeKey(panel.Key)
		panel.Label = strings.TrimSpace(panel.Label)
		panel.Summary = strings.TrimSpace(panel.Summary)
		if panel.Key == "" || panel.Label == "" {
			continue
		}
		out = append(out, panel)
	}
	return out
}

func hasAction(actions []Action, key string) bool {
	key = normalizeKey(key)
	for _, action := range actions {
		if normalizeKey(action.Key) == key {
			return true
		}
	}
	return false
}

func normalizeKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}
