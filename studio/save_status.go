package studio

import (
	"m31labs.dev/gosx"
	gosxstudio "m31labs.dev/gosx-studio"
)

type SaveStatusOptions struct {
	Class           string
	StateClass      string
	DetailClass     string
	LastSavedClass  string
	DirtyCountClass string
	StateLabel      string
	DetailLabel     string
	LastSavedLabel  string
}

func RenderSaveStatus(options SaveStatusOptions) gosx.Node {
	className := gosxstudio.FirstNonEmpty(options.Class, "gosx-studio-save-status")
	stateClass := gosxstudio.FirstNonEmpty(options.StateClass, "gosx-studio-save-status__state")
	detailClass := gosxstudio.FirstNonEmpty(options.DetailClass, "gosx-studio-save-status__detail")
	lastSavedClass := gosxstudio.FirstNonEmpty(options.LastSavedClass, "gosx-studio-save-status__last-saved")
	dirtyCountClass := gosxstudio.FirstNonEmpty(options.DirtyCountClass, "gosx-studio-save-status__dirty-count")
	stateLabel := gosxstudio.FirstNonEmpty(options.StateLabel, "Saved")
	detailLabel := gosxstudio.FirstNonEmpty(options.DetailLabel, "Ready")
	lastSavedLabel := gosxstudio.FirstNonEmpty(options.LastSavedLabel, "Not saved this session")

	return gosx.El("div", gosx.Attrs(
		gosx.Attr("class", className),
		gosx.Attr("data-gosx-studio-save-status", "true"),
	),
		gosx.El("output", gosx.Attrs(
			gosx.Attr("class", stateClass),
			gosx.Attr("data-gosx-studio-save-state", "true"),
			gosx.Attr("aria-live", "polite"),
		), gosx.Text(stateLabel)),
		gosx.El("span", gosx.Attrs(
			gosx.Attr("class", detailClass),
			gosx.Attr("data-gosx-studio-save-detail", "true"),
		), gosx.Text(detailLabel)),
		gosx.El("output", gosx.Attrs(
			gosx.Attr("class", dirtyCountClass),
			gosx.Attr("data-gosx-studio-dirty-count", "true"),
			gosx.Attr("hidden", "hidden"),
		), gosx.Text("0 changes")),
		gosx.El("time", gosx.Attrs(
			gosx.Attr("class", lastSavedClass),
			gosx.Attr("data-gosx-studio-last-saved", "true"),
			gosx.Attr("data-gosx-studio-last-saved-empty", lastSavedLabel),
		), gosx.Text(lastSavedLabel)),
	)
}
