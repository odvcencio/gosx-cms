package studio

import (
	_ "embed"

	"github.com/odvcencio/gosx"
)

//go:embed assets/command_palette.js
var commandPaletteRuntime string

//go:embed assets/state_runtime.js
var stateRuntime string

//go:embed assets/site_canvas_runtime.js
var siteCanvasRuntime string

func CommandPaletteScript() string {
	return commandPaletteRuntime
}

func RenderCommandPaletteScript() gosx.Node {
	return gosx.RawHTML(`<script data-gosx-studio-command-runtime="true">` + commandPaletteRuntime + `</script>`)
}

func StateRuntimeScript() string {
	return stateRuntime
}

func RenderStudioStateScript() gosx.Node {
	return gosx.RawHTML(`<script data-gosx-studio-state-runtime="true">` + stateRuntime + `</script>`)
}

func SiteCanvasScript() string {
	return siteCanvasRuntime
}

func RenderSiteCanvasScript() gosx.Node {
	return gosx.RawHTML(`<script data-gosx-studio-site-canvas-runtime="true">` + siteCanvasRuntime + `</script>`)
}
