package studio

import (
	"strings"
	"testing"

	"github.com/odvcencio/gosx"
	cmsflows "github.com/odvcencio/gosx-cms/flows"
)

func TestRenderFlowEditorKeepsAuthoringContracts(t *testing.T) {
	html := gosx.RenderHTML(RenderFlowEditor([]cmsflows.StudioFlow{
		{
			Key:         "schedule-tour",
			Label:       "Schedule tour",
			Summary:     "2 steps / 1 action / 1 field",
			StatusLabel: "Ready",
			Route:       "/schedule-tour",
			HasRoute:    true,
			PrimaryAction: cmsflows.StudioAction{
				HandlerRef: "tour.submit",
			},
			Steps: []cmsflows.StudioStep{
				{Key: "request-info", Label: "Request info"},
			},
		},
	}, FlowEditorOptions{PublishAction: "/admin/editor/flows/publish"}))
	for _, check := range []string{
		`class="list-card studio-flow-summary is-selected"`,
		`id="studio-flow-tab-schedule-tour"`,
		`data-studio-flow-library="true"`,
		`data-studio-flow-card="schedule-tour"`,
		`data-studio-flow-route="/schedule-tour"`,
		`aria-controls="studio-flow-editor-schedule-tour"`,
		`data-studio-flow-dirty-badge="true"`,
		`data-studio-flow-fields="true"`,
		`data-studio-flow-editor="schedule-tour"`,
		`id="studio-flow-handler-schedule-tour"`,
		`name="flowScheduleTourHandlerRef"`,
		`value="tour.submit" data-studio-initial-value="tour.submit"`,
		`id="studio-flow-step-schedule-tour-request-info-label"`,
		`name="flowScheduleTourStepRequestInfoLabel"`,
		`data-studio-preview-flow="/schedule-tour"`,
		`formaction="/admin/editor/flows/publish" name="flowKey" value="schedule-tour"`,
	} {
		if !strings.Contains(html, check) {
			t.Fatalf("expected %q in flow editor markup: %s", check, html)
		}
	}
}

func TestFlowEditorInputNameHelpersUseStablePrefixes(t *testing.T) {
	if got := FlowHandlerRefInputName("schedule-tour"); got != "flowScheduleTourHandlerRef" {
		t.Fatalf("unexpected handler input name: %q", got)
	}
	if got := FlowStepLabelInputName("schedule-tour", "request.info"); got != "flowScheduleTourStepRequestInfoLabel" {
		t.Fatalf("unexpected step input name: %q", got)
	}
}
