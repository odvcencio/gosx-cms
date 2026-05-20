package studio

import (
	_ "embed"

	"github.com/odvcencio/gosx"
)

//go:embed assets/command_palette.js
var commandPaletteRuntime string

func CommandPaletteScript() string {
	return commandPaletteRuntime
}

func RenderCommandPaletteScript() gosx.Node {
	return gosx.RawHTML(`<script data-gosx-studio-command-runtime="true">` + commandPaletteRuntime + `</script>`)
}
