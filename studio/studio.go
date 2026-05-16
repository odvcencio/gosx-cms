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
	Navigation    []Section
	Metrics       []Metric
	Left          []Panel
	Right         []Panel
	Actions       []Action
	Extras        map[string]any
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

type Section struct {
	Key     string
	Label   string
	Summary string
	Actions []Action
}

type Metric struct {
	Key   string
	Label string
	Value any
}

type BlockSummary struct {
	Key        string
	Label      string
	Summary    string
	Kind       string
	Preview    string
	DefaultOn  bool
	Locked     bool
	Repeatable bool
	Icon       string
	FieldCount int
}

type MediaSummary struct {
	ID          string
	URL         string
	Alt         string
	Filename    string
	ContentType string
	Size        int64
	Archived    bool
}

type RevisionSummary struct {
	ID            string
	ResourceKind  string
	ResourceID    string
	ResourceTitle string
	Action        string
	Summary       string
	Created       string
}

type Shell struct {
	Title         string
	PreviewURL    string
	SaveAction    string
	RestoreAction string
	BlockCatalog  []blockstudio.Definition
	Blocks        []BlockSummary
	BlockCount    int
	Media         []media.Asset
	MediaLibrary  []MediaSummary
	MediaCount    int
	HasMedia      bool
	Revisions     []lifecycle.Revision
	RevisionLog   []RevisionSummary
	RevisionCount int
	HasRevisions  bool
	Navigation    []Section
	Metrics       []Metric
	Left          []Panel
	Right         []Panel
	Actions       []Action
	Extras        map[string]any
}

func View(shell Shell) map[string]any {
	view := map[string]any{
		"title":         shell.Title,
		"previewURL":    shell.PreviewURL,
		"saveAction":    shell.SaveAction,
		"restoreAction": shell.RestoreAction,
		"blockCatalog":  blockViews(shell.Blocks),
		"blockCount":    shell.BlockCount,
		"media":         mediaViews(shell.MediaLibrary),
		"mediaCount":    shell.MediaCount,
		"hasMedia":      shell.HasMedia,
		"revisions":     revisionViews(shell.RevisionLog),
		"revisionCount": shell.RevisionCount,
		"hasRevisions":  shell.HasRevisions,
		"navigation":    sectionViews(shell.Navigation),
		"metrics":       metricViews(shell.Metrics),
		"actions":       actionViews(shell.Actions),
		"leftPanels":    panelViews(shell.Left),
		"rightPanels":   panelViews(shell.Right),
	}
	for key, value := range shell.Extras {
		if _, exists := view[key]; !exists {
			view[key] = value
		}
	}
	if shell.Extras != nil {
		view["extras"] = cloneExtras(shell.Extras)
	}
	return view
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
	navigation := normalizeSections(options.Navigation)
	metrics := normalizeMetrics(options.Metrics)
	blockCatalog := append([]blockstudio.Definition(nil), options.BlockCatalog...)
	mediaAssets := media.CloneAssets(options.Media)
	revisions := lifecycle.CloneRevisions(options.Revisions)
	return Shell{
		Title:         title,
		PreviewURL:    previewURL,
		SaveAction:    strings.TrimSpace(options.SaveAction),
		RestoreAction: strings.TrimSpace(options.RestoreAction),
		BlockCatalog:  blockCatalog,
		Blocks:        blockSummaries(blockCatalog),
		BlockCount:    len(blockCatalog),
		Media:         mediaAssets,
		MediaLibrary:  mediaSummaries(mediaAssets),
		MediaCount:    len(mediaAssets),
		HasMedia:      len(mediaAssets) > 0,
		Revisions:     revisions,
		RevisionLog:   revisionSummaries(revisions),
		RevisionCount: len(revisions),
		HasRevisions:  len(revisions) > 0,
		Navigation:    navigation,
		Metrics:       metrics,
		Left:          normalizePanels(options.Left),
		Right:         normalizePanels(options.Right),
		Actions:       actions,
		Extras:        cloneExtras(options.Extras),
	}
}

func LinkAction(key, label, href string) Action {
	return Action{Key: key, Label: label, Href: href}
}

func PrimaryAction(key, label, href string) Action {
	return Action{Key: key, Label: label, Href: href, Primary: true}
}

func NewPanel(key, label, summary string, children ...gosx.Node) Panel {
	return Panel{Key: key, Label: label, Summary: summary, Children: children}
}

func NewSection(key, label string, actions ...Action) Section {
	return Section{Key: key, Label: label, Actions: actions}
}

func NewMetric(key, label string, value any) Metric {
	return Metric{Key: key, Label: label, Value: value}
}

func actionViews(actions []Action) []map[string]any {
	out := make([]map[string]any, 0, len(actions))
	for _, action := range actions {
		className := "button button--secondary"
		if action.Primary {
			className = "button button--primary"
		}
		out = append(out, map[string]any{
			"key":     action.Key,
			"label":   action.Label,
			"href":    action.Href,
			"primary": action.Primary,
			"class":   className,
		})
	}
	return out
}

func sectionViews(sections []Section) []map[string]any {
	out := make([]map[string]any, 0, len(sections))
	for _, section := range sections {
		out = append(out, map[string]any{
			"key":        section.Key,
			"label":      section.Label,
			"summary":    section.Summary,
			"hasSummary": section.Summary != "",
			"actions":    actionViews(section.Actions),
		})
	}
	return out
}

func metricViews(metrics []Metric) []map[string]any {
	out := make([]map[string]any, 0, len(metrics))
	for _, metric := range metrics {
		out = append(out, map[string]any{
			"key":   metric.Key,
			"label": metric.Label,
			"value": metric.Value,
		})
	}
	return out
}

func panelViews(panels []Panel) []map[string]any {
	out := make([]map[string]any, 0, len(panels))
	for _, panel := range panels {
		out = append(out, map[string]any{
			"key":        panel.Key,
			"label":      panel.Label,
			"summary":    panel.Summary,
			"hasSummary": panel.Summary != "",
		})
	}
	return out
}

func blockViews(blocks []BlockSummary) []map[string]any {
	out := make([]map[string]any, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, map[string]any{
			"key":        block.Key,
			"label":      block.Label,
			"summary":    block.Summary,
			"kind":       block.Kind,
			"preview":    block.Preview,
			"defaultOn":  block.DefaultOn,
			"locked":     block.Locked,
			"repeatable": block.Repeatable,
			"icon":       block.Icon,
			"fieldCount": block.FieldCount,
		})
	}
	return out
}

func mediaViews(media []MediaSummary) []map[string]any {
	out := make([]map[string]any, 0, len(media))
	for _, asset := range media {
		out = append(out, map[string]any{
			"id":          asset.ID,
			"url":         asset.URL,
			"alt":         asset.Alt,
			"filename":    asset.Filename,
			"contentType": asset.ContentType,
			"size":        asset.Size,
			"archived":    asset.Archived,
		})
	}
	return out
}

func revisionViews(revisions []RevisionSummary) []map[string]any {
	out := make([]map[string]any, 0, len(revisions))
	for _, revision := range revisions {
		out = append(out, map[string]any{
			"id":            revision.ID,
			"resourceKind":  revision.ResourceKind,
			"resourceID":    revision.ResourceID,
			"resourceTitle": revision.ResourceTitle,
			"action":        revision.Action,
			"summary":       revision.Summary,
			"created":       revision.Created,
		})
	}
	return out
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

func normalizeSections(sections []Section) []Section {
	out := make([]Section, 0, len(sections))
	for _, section := range sections {
		section.Key = normalizeKey(section.Key)
		section.Label = strings.TrimSpace(section.Label)
		section.Summary = strings.TrimSpace(section.Summary)
		section.Actions = normalizeActions(section.Actions, "")
		if section.Key == "" || section.Label == "" {
			continue
		}
		out = append(out, section)
	}
	return out
}

func normalizeMetrics(metrics []Metric) []Metric {
	out := make([]Metric, 0, len(metrics))
	for _, metric := range metrics {
		metric.Key = normalizeKey(metric.Key)
		metric.Label = strings.TrimSpace(metric.Label)
		if metric.Key == "" || metric.Label == "" {
			continue
		}
		out = append(out, metric)
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

func blockSummaries(blocks []blockstudio.Definition) []BlockSummary {
	out := make([]BlockSummary, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, BlockSummary{
			Key:        block.Key,
			Label:      block.Label,
			Summary:    block.Summary,
			Kind:       block.Kind,
			Preview:    block.Preview,
			DefaultOn:  block.DefaultOn,
			Locked:     block.Locked,
			Repeatable: block.Repeatable,
			Icon:       block.Icon,
			FieldCount: len(block.Fields),
		})
	}
	return out
}

func mediaSummaries(assets []media.Asset) []MediaSummary {
	out := make([]MediaSummary, 0, len(assets))
	for _, asset := range assets {
		out = append(out, MediaSummary{
			ID:          asset.ID,
			URL:         asset.URL,
			Alt:         asset.Alt,
			Filename:    asset.Filename,
			ContentType: asset.ContentType,
			Size:        asset.Size,
			Archived:    asset.ArchivedAt != nil,
		})
	}
	return out
}

func revisionSummaries(revisions []lifecycle.Revision) []RevisionSummary {
	out := make([]RevisionSummary, 0, len(revisions))
	for _, revision := range revisions {
		created := ""
		if !revision.Created.IsZero() {
			created = revision.Created.Format("2006-01-02T15:04:05Z07:00")
		}
		out = append(out, RevisionSummary{
			ID:            revision.ID,
			ResourceKind:  revision.ResourceKind,
			ResourceID:    revision.ResourceID,
			ResourceTitle: revision.ResourceTitle,
			Action:        revision.Action,
			Summary:       revision.Summary,
			Created:       created,
		})
	}
	return out
}

func cloneExtras(extras map[string]any) map[string]any {
	if extras == nil {
		return nil
	}
	out := make(map[string]any, len(extras))
	for key, value := range extras {
		out[key] = value
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
