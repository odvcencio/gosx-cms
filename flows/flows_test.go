package flows

import "testing"

func TestCatalogNormalizesAndDedupes(t *testing.T) {
	catalog := Catalog(
		Definition{Key: "Schedule Request", Label: "Schedule request", Steps: []Step{{Label: "Request"}}, Actions: []Action{{Key: "submit", HandlerRef: "schedule.submit"}}},
		Definition{Key: "schedule-request", Label: "Duplicate", Steps: []Step{{Key: "request"}}},
	)
	if len(catalog) != 1 {
		t.Fatalf("expected one deduped flow, got %#v", catalog)
	}
	flow := catalog[0]
	if flow.Key != "schedule-request" || flow.Steps[0].Key != "request" || flow.Steps[0].Blocks.Version != 1 {
		t.Fatalf("unexpected normalized flow: %#v", flow)
	}
}

func TestValidateRequiresStepsAndHandlers(t *testing.T) {
	errs := Validate(Definition{Key: "contact", Actions: []Action{{Key: "submit"}}})
	if errs["steps"] == "" || errs["actions.submit"] == "" {
		t.Fatalf("expected validation errors, got %#v", errs)
	}
}

func TestFind(t *testing.T) {
	catalog := Catalog(Contact("contact.submit"))
	flow, ok := Find(catalog, "contact")
	if !ok || flow.Label != "Contact" {
		t.Fatalf("expected contact flow, got %#v %v", flow, ok)
	}
}

func TestPresetFlows(t *testing.T) {
	for _, flow := range []Definition{
		Contact("contact.submit"),
		ScheduleRequest("schedule.submit"),
		Enrollment("enrollment.submit"),
	} {
		if errs := Validate(flow); len(errs) != 0 {
			t.Fatalf("expected preset flow to validate: %#v %#v", flow, errs)
		}
		if len(flow.Actions) != 1 || len(flow.Actions[0].Fields) == 0 {
			t.Fatalf("expected preset flow fields: %#v", flow)
		}
	}
}
