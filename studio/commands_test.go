package studio

import (
	"strings"
	"testing"

	"github.com/odvcencio/gosx"
)

func TestRenderCommandPalette(t *testing.T) {
	html := gosx.RenderHTML(RenderCommandPalette(CommandPaletteOptions{
		Class: "studio-command-palette",
		Commands: []Command{
			{Kind: CommandMode, Key: "mode-structure", Label: "Structure", Group: "Modes", Target: "structure", Shortcut: "1"},
			{Kind: CommandInsert, Label: "Hero", Target: "hero", Summary: "Add hero block", Primary: true, Keywords: []string{"home"}},
			{Kind: CommandLink, Key: "media", Label: "Media", Href: "/admin/media"},
			{Label: ""},
		},
	}))
	if !strings.Contains(html, `data-studio-command-palette="true"`) || !strings.Contains(html, `data-studio-command-open="true"`) || !strings.Contains(html, `data-studio-command-search="true"`) {
		t.Fatalf("expected command palette shell hooks, got: %s", html)
	}
	if !strings.Contains(html, `data-studio-command-kind="mode"`) || !strings.Contains(html, `data-studio-command-target="structure"`) {
		t.Fatalf("expected mode command hooks, got: %s", html)
	}
	if !strings.Contains(html, `studio-command-item--primary`) || !strings.Contains(html, `data-studio-command-kind="insert"`) {
		t.Fatalf("expected primary insert command, got: %s", html)
	}
	if !strings.Contains(html, `data-studio-command-href="/admin/media"`) {
		t.Fatalf("expected link command href hook, got: %s", html)
	}
	if strings.Contains(html, `data-studio-command-kind=""`) {
		t.Fatalf("expected invalid commands to be skipped: %s", html)
	}
}

func TestNormalizeCommandsDedupeAndDefaults(t *testing.T) {
	commands := normalizeCommands([]Command{
		{Kind: CommandMode, Label: "Style", Target: "style"},
		{Kind: CommandMode, Label: "Style duplicate", Target: "style"},
		{Kind: "unknown", Key: "media", Label: "Media", Href: "/admin/media"},
	})
	if len(commands) != 2 {
		t.Fatalf("expected generated duplicate key to dedupe, got %#v", commands)
	}
	if commands[0].Key != "mode-style" || commands[0].Group != "Mode" {
		t.Fatalf("unexpected generated command: %#v", commands[0])
	}
	if commands[1].Kind != CommandLink || commands[1].Group != "Open" {
		t.Fatalf("unexpected unknown command normalization: %#v", commands[1])
	}
}
