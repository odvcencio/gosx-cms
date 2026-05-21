package studio

import (
	"fmt"
	"strings"

	"github.com/odvcencio/gosx"
)

type MetricCardsOptions struct {
	Class     string
	CardClass string
}

type FlowStorageSummary struct {
	Kind          string
	Class         string
	Label         string
	StatusLabel   string
	Detail        string
	Path          string
	DocumentCount int
}

type FlowStoragePanelOptions struct {
	Class       string
	HeadClass   string
	KickerClass string
	MetaClass   string
	Kicker      string
}

type LifecycleActivityItem struct {
	Summary string
	Action  string
	Created string
}

type LifecycleControlsPanelOptions struct {
	Class             string
	FieldClass        string
	ButtonRowClass    string
	ActivityClass     string
	ActivityCardClass string
	Title             string
	PublishAtName     string
	ApproveAction     string
	ScheduleAction    string
	ProcessDueAction  string
	ApproveLabel      string
	ScheduleLabel     string
	ProcessDueLabel   string
	Activity          []LifecycleActivityItem
}

type StudioMapPanelOptions struct {
	Class       string
	GridClass   string
	PanelClass  string
	CardClass   string
	Title       string
	PanelsTitle string
}

type WorkbenchPreviewPanelOptions struct {
	Class        string
	Title        string
	URL          string
	IFrameTitle  string
	OpenLabel    string
	OpenClass    string
	FrameClass   string
	ToolbarClass string
	MetaClass    string
	StatusClass  string
}

type WorkbenchSummaryToolbarOptions struct {
	Class              string
	ActionsClass       string
	Title              string
	Summary            string
	SaveButtonClass    string
	SaveButtonLabel    string
	CommandPaletteNode gosx.Node
	SaveStatusNode     gosx.Node
	Controls           []gosx.Node
	Actions            []gosx.Node
}

func RenderMetricCards(metrics []Metric, options MetricCardsOptions) gosx.Node {
	cards := make([]gosx.Node, 0, len(metrics))
	for _, metric := range normalizeMetrics(metrics) {
		cards = append(cards, gosx.El("article", gosx.Attrs(
			gosx.Attr("class", firstNonEmpty(options.CardClass, "stat-card")),
			gosx.Attr("data-studio-metric", metric.Key),
		),
			gosx.El("span", nil, gosx.Text(metric.Label)),
			gosx.El("strong", nil, gosx.Text(fmtAny(metric.Value))),
		))
	}
	return gosx.El("section", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.Class, "stat-grid"))), gosx.Fragment(cards...))
}

func RenderFlowStoragePanel(summary FlowStorageSummary, options FlowStoragePanelOptions) gosx.Node {
	children := []gosx.Node{
		gosx.El("div", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.HeadClass, "studio-storage__head"))),
			gosx.El("div", nil,
				gosx.El("p", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.KickerClass, "studio-storage__kicker"))), gosx.Text(firstNonEmpty(options.Kicker, "Storage"))),
				gosx.El("h2", nil, gosx.Text(summary.Label)),
			),
			gosx.El("output", nil, gosx.Text(summary.StatusLabel)),
		),
		gosx.El("p", nil, gosx.Text(summary.Detail)),
	}
	meta := []gosx.Node{
		gosx.El("span", nil,
			gosx.El("strong", nil, gosx.Text(fmtAny(summary.DocumentCount))),
			gosx.Text(" flow documents"),
		),
	}
	if path := strings.TrimSpace(summary.Path); path != "" {
		meta = append(meta, gosx.El("code", nil, gosx.Text(path)))
	}
	children = append(children, gosx.El("div", gosx.Attrs(
		gosx.Attr("class", firstNonEmpty(options.MetaClass, "studio-storage__meta")),
		gosx.Attr("aria-label", "Flow storage summary"),
	), gosx.Fragment(meta...)))
	return gosx.El("section", gosx.Attrs(
		gosx.Attr("class", firstNonEmpty(options.Class, summary.Class, "studio-panel studio-storage")),
		gosx.Attr("data-studio-flow-storage", summary.Kind),
	), gosx.Fragment(children...))
}

func RenderLifecycleControlsPanel(state LifecycleReviewState, options LifecycleControlsPanelOptions) gosx.Node {
	if strings.TrimSpace(options.ApproveAction) == "" && strings.TrimSpace(options.ScheduleAction) == "" && strings.TrimSpace(options.ProcessDueAction) == "" {
		return gosx.Fragment()
	}
	buttons := []gosx.Node{}
	if strings.TrimSpace(options.ApproveAction) != "" {
		buttons = append(buttons, gosx.El("button", gosx.Attrs(
			gosx.Attr("class", "button button--secondary"),
			gosx.Attr("type", "submit"),
			gosx.Attr("formaction", options.ApproveAction),
		), gosx.Text(firstNonEmpty(options.ApproveLabel, "Approve publish"))))
	}
	if strings.TrimSpace(options.ScheduleAction) != "" {
		buttons = append(buttons, gosx.El("button", gosx.Attrs(
			gosx.Attr("class", "button button--primary"),
			gosx.Attr("type", "submit"),
			gosx.Attr("formaction", options.ScheduleAction),
		), gosx.Text(firstNonEmpty(options.ScheduleLabel, "Schedule publish"))))
	}
	if strings.TrimSpace(options.ProcessDueAction) != "" {
		buttons = append(buttons, gosx.El("button", gosx.Attrs(
			gosx.Attr("class", "button button--secondary"),
			gosx.Attr("type", "submit"),
			gosx.Attr("formaction", options.ProcessDueAction),
		), gosx.Text(firstNonEmpty(options.ProcessDueLabel, "Run due publishes"))))
	}
	children := []gosx.Node{
		gosx.El("h2", nil, gosx.Text(firstNonEmpty(options.Title, "Lifecycle"))),
		gosx.El("label", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.FieldClass, "field"))),
			gosx.El("span", nil, gosx.Text("Publish time")),
			gosx.El("input", gosx.Attrs(
				gosx.Attr("type", "datetime-local"),
				gosx.Attr("name", firstNonEmpty(options.PublishAtName, "publishAt")),
				gosx.Attr("value", LifecycleScheduleInputValue(state, nil)),
			)),
			gosx.El("small", nil, gosx.Text(LifecycleScheduleHelp(state, nil))),
		),
		gosx.El("div", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.ButtonRowClass, "button-row"))), gosx.Fragment(buttons...)),
	}
	if len(options.Activity) > 0 {
		events := make([]gosx.Node, 0, len(options.Activity))
		for _, event := range options.Activity {
			events = append(events, gosx.El("article", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.ActivityCardClass, "list-card"))),
				gosx.El("strong", nil, gosx.Text(event.Summary)),
				gosx.El("span", nil, gosx.Text(strings.TrimSpace(event.Action+" / "+event.Created))),
			))
		}
		children = append(children, gosx.El("div", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.ActivityClass, "stack"))), gosx.Fragment(events...)))
	}
	return gosx.El("section", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.Class, "studio-panel"))), gosx.Fragment(children...))
}

func RenderStudioMapPanel(shell Shell, options StudioMapPanelOptions) gosx.Node {
	navItems := make([]gosx.Node, 0, len(shell.Navigation))
	for _, section := range shell.Navigation {
		navItems = append(navItems, gosx.El("article", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.CardClass, "list-card"))),
			gosx.El("strong", nil, gosx.Text(section.Label)),
			gosx.El("span", nil, gosx.Text(section.Summary)),
		))
	}
	panelItems := make([]gosx.Node, 0, len(shell.Left)+len(shell.Right))
	for _, panel := range append(append([]Panel{}, shell.Left...), shell.Right...) {
		panelItems = append(panelItems, gosx.El("article", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.CardClass, "list-card"))),
			gosx.El("strong", nil, gosx.Text(panel.Label)),
			gosx.El("span", nil, gosx.Text(panel.Summary)),
		))
	}
	return gosx.El("div", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.GridClass, options.Class, "studio-grid"))),
		gosx.El("section", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.PanelClass, "studio-panel"))),
			gosx.El("h2", nil, gosx.Text(firstNonEmpty(options.Title, "Studio map"))),
			gosx.Fragment(navItems...),
		),
		gosx.El("section", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.PanelClass, "studio-panel"))),
			gosx.El("h2", nil, gosx.Text(firstNonEmpty(options.PanelsTitle, "Panels"))),
			gosx.Fragment(panelItems...),
		),
	)
}

func RenderWorkbenchPreviewPanel(shell Shell, options WorkbenchPreviewPanelOptions) gosx.Node {
	url := firstNonEmpty(options.URL, shell.PreviewURL, "/")
	return gosx.El("section", gosx.Attrs(gosx.Attr("class", firstNonEmpty(options.Class, "studio-panel studio-preview-panel"))),
		gosx.El("h2", nil, gosx.Text(firstNonEmpty(options.Title, "Preview"))),
		RenderPreviewFrame(PreviewFrameOptions{
			ShellClass:   "studio-preview-frame-shell",
			ToolbarClass: firstNonEmpty(options.ToolbarClass, "studio-preview-toolbar"),
			MetaClass:    firstNonEmpty(options.MetaClass, "studio-preview-toolbar__flow"),
			FrameClass:   firstNonEmpty(options.FrameClass, "preview-frame"),
			StatusClass:  firstNonEmpty(options.StatusClass, "studio-preview-toolbar__status"),
			OpenClass:    firstNonEmpty(options.OpenClass, "button button--secondary"),
			Kicker:       "Previewing",
			Title:        "Public site",
			URL:          url,
			IFrameTitle:  firstNonEmpty(options.IFrameTitle, "Public site preview"),
			StatusLabel:  "Ready",
			OpenLabel:    firstNonEmpty(options.OpenLabel, "Open route"),
			OpenNewTab:   true,
			DynamicTitle: true,
			DynamicRoute: true,
		}),
	)
}

func RenderWorkbenchSummaryToolbar(shell Shell, options WorkbenchSummaryToolbarOptions) gosx.Node {
	actions := []gosx.Node{}
	if !isZeroNode(options.CommandPaletteNode) {
		actions = append(actions, options.CommandPaletteNode)
	}
	if !isZeroNode(options.SaveStatusNode) {
		actions = append(actions, options.SaveStatusNode)
	}
	actions = append(actions, options.Actions...)
	if firstNonEmpty(options.SaveButtonLabel, "Save checkpoint") != "" {
		actions = append(actions, gosx.El("button", gosx.Attrs(
			gosx.Attr("class", firstNonEmpty(options.SaveButtonClass, "button button--primary")),
			gosx.Attr("type", "submit"),
			gosx.Attr("data-gosx-studio-save-button", "true"),
		), gosx.Text(firstNonEmpty(options.SaveButtonLabel, "Save checkpoint"))))
	}
	return RenderStudioToolbar(StudioToolbarOptions{
		Class:        firstNonEmpty(options.Class, "studio-toolbar"),
		ActionsClass: firstNonEmpty(options.ActionsClass, "studio-toolbar__actions"),
		Title:        firstNonEmpty(options.Title, "Authoring surfaces"),
		Summary:      firstNonEmpty(options.Summary, fmt.Sprintf("%d reusable blocks", shell.BlockCount)),
		Controls:     options.Controls,
		Actions:      actions,
	})
}

func MetricValue(metrics []Metric, key string) string {
	for _, metric := range normalizeMetrics(metrics) {
		if metric.Key == normalizeKey(key) {
			return fmtAny(metric.Value)
		}
	}
	return "0"
}
