package studio

import "strings"

type CommandBlock struct {
	Key          string
	Label        string
	Summary      string
	Target       string
	PreviewTitle string
}

type CommandFlow struct {
	Key            string
	Label          string
	Description    string
	Route          string
	EmbedTarget    string
	HasRoute       bool
	HasEmbedTarget bool
}

type StudioCommandOptions struct {
	Shell        Shell
	Blocks       []CommandBlock
	Flows        []CommandFlow
	Extra        []Command
	SaveLabel    string
	SaveSummary  string
	SaveKeywords []string
}

func StudioCommands(options StudioCommandOptions) []Command {
	shell := options.Shell
	commands := []Command{
		{
			Kind:     CommandSubmit,
			Key:      "save",
			Label:    firstNonEmpty(options.SaveLabel, "Save changes"),
			Summary:  firstNonEmpty(options.SaveSummary, "Persist the current draft."),
			Group:    "Save",
			Target:   "save",
			Shortcut: "Ctrl S",
			Primary:  true,
			Keywords: firstKeywords(options.SaveKeywords, []string{"publish", "draft", "settings"}),
		},
		{Kind: CommandHistory, Key: "undo", Label: "Undo", Summary: "Revert the last in-browser editor change.", Group: "History", Target: "undo", Shortcut: "Ctrl Z", Keywords: []string{"back", "revert"}},
		{Kind: CommandHistory, Key: "redo", Label: "Redo", Summary: "Reapply the last undone editor change.", Group: "History", Target: "redo", Shortcut: "Ctrl Shift Z", Keywords: []string{"forward", "restore"}},
		{Kind: CommandToggle, Key: "toggle-layers", Label: "Toggle layers rail", Summary: "Show or hide the page and layer navigator.", Group: "View", Target: "left", Shortcut: "L", Keywords: []string{"sidebar", "left"}},
		{Kind: CommandToggle, Key: "toggle-inspector", Label: "Toggle inspector rail", Summary: "Show or hide properties and style controls.", Group: "View", Target: "right", Shortcut: "I", Keywords: []string{"sidebar", "properties"}},
		{Kind: CommandToggle, Key: "toggle-activity", Label: "Toggle activity rail", Summary: "Show readiness checks, proposals, and review activity.", Group: "View", Target: "activity", Shortcut: "A", Keywords: []string{"comments", "proposals"}},
		{Kind: CommandToggle, Key: "toggle-focus", Label: "Focus canvas", Summary: "Give the live preview more room.", Group: "View", Target: "focus", Shortcut: "F", Keywords: []string{"preview", "canvas"}},
	}
	for _, mode := range shell.Modes {
		commands = append(commands, Command{
			Kind:     CommandMode,
			Key:      "mode-" + mode.Key,
			Label:    mode.Label + " mode",
			Summary:  "Switch the editor to " + strings.ToLower(mode.Label) + ".",
			Group:    "Modes",
			Target:   mode.Key,
			Keywords: []string{"panel", "inspector"},
		})
	}
	for _, viewport := range shell.Viewports {
		commands = append(commands, Command{
			Kind:     CommandViewport,
			Key:      "viewport-" + viewport.Key,
			Label:    viewport.Label + " viewport",
			Summary:  "Preview the page at " + strings.ToLower(viewport.Label) + " size.",
			Group:    "Preview",
			Target:   viewport.Key,
			Keywords: []string{"responsive", "breakpoint"},
		})
	}
	for _, zoom := range DefaultZoomLevels(shell.Canvas.Zoom) {
		commands = append(commands, Command{
			Kind:   CommandZoom,
			Key:    "zoom-" + zoom.Key,
			Label:  "Zoom " + zoom.Label,
			Group:  "Preview",
			Target: zoom.Key,
		})
	}
	for _, action := range DefaultSelectionCommands() {
		commands = append(commands, Command{
			Kind:     CommandSelectionAction,
			Key:      "selection-" + action.Key,
			Label:    action.Label + " selected block",
			Summary:  "Run this on the current block selection.",
			Group:    "Selection",
			Target:   action.Key,
			Keywords: []string{"block", "canvas"},
		})
	}
	for _, block := range options.Blocks {
		label := strings.TrimSpace(block.Label)
		target := normalizeKey(firstNonEmpty(block.Target, block.Key))
		if label == "" || target == "" {
			continue
		}
		commands = append(commands, Command{
			Kind:     CommandInsert,
			Key:      "insert-" + target,
			Label:    "Add " + label,
			Summary:  "Insert or enable this homepage block.",
			Group:    "Blocks",
			Target:   target,
			Keywords: []string{block.PreviewTitle, block.Summary},
		})
	}
	commands = append(commands, options.Extra...)
	for _, flow := range options.Flows {
		flow.Key = normalizeKey(flow.Key)
		flow.Label = strings.TrimSpace(flow.Label)
		flow.Description = strings.TrimSpace(flow.Description)
		flow.Route = strings.TrimSpace(flow.Route)
		flow.EmbedTarget = normalizeKey(flow.EmbedTarget)
		if flow.Key == "" || flow.Label == "" {
			continue
		}
		if flow.HasRoute && flow.Route != "" {
			commands = append(commands, Command{
				Kind:     CommandLink,
				Key:      "open-flow-" + flow.Key,
				Label:    "Open " + flow.Label,
				Summary:  flow.Description,
				Group:    "Flows",
				Href:     flow.Route,
				Keywords: []string{"flow", "form", "behavior"},
			})
		}
		if flow.HasEmbedTarget && flow.EmbedTarget != "" {
			commands = append(commands, Command{
				Kind:     CommandInsert,
				Key:      "insert-flow-" + flow.Key,
				Label:    "Add " + flow.Label + " flow",
				Summary:  "Embed this behavior on the current page.",
				Group:    "Flows",
				Target:   flow.EmbedTarget,
				Keywords: []string{"flow", "form", "behavior"},
			})
		}
	}
	return normalizeCommands(commands)
}

func firstKeywords(values, fallback []string) []string {
	if len(values) == 0 {
		return append([]string{}, fallback...)
	}
	out := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = append(out, value)
		}
	}
	if len(out) == 0 {
		return append([]string{}, fallback...)
	}
	return out
}
