package flows

import (
	"fmt"
	"strings"

	"github.com/odvcencio/gosx-admin/blockstudio"
	"github.com/odvcencio/gosx-admin/workbench"
)

type Definition struct {
	Key         string
	Label       string
	Description string
	Steps       []Step
	Actions     []Action
}

type Step struct {
	Key    string
	Label  string
	Blocks blockstudio.Document
}

type Action struct {
	Key        string
	Label      string
	HandlerRef string
	Fields     []workbench.Field
}

type ValidationErrors map[string]string

func Normalize(definition Definition) Definition {
	definition.Key = normalizeKey(definition.Key)
	definition.Label = strings.TrimSpace(definition.Label)
	definition.Description = strings.TrimSpace(definition.Description)
	if definition.Label == "" {
		definition.Label = definition.Key
	}
	steps := make([]Step, 0, len(definition.Steps))
	seenSteps := map[string]bool{}
	for index, step := range definition.Steps {
		step.Key = normalizeKey(firstNonEmpty(step.Key, step.Label))
		if step.Key == "" {
			step.Key = "step"
		}
		if seenSteps[step.Key] {
			step.Key = fmt.Sprintf("%s-%d", step.Key, index+1)
		}
		seenSteps[step.Key] = true
		step.Label = firstNonEmpty(step.Label, step.Key)
		if step.Blocks.Version <= 0 {
			step.Blocks.Version = 1
		}
		steps = append(steps, step)
	}
	definition.Steps = steps
	actions := make([]Action, 0, len(definition.Actions))
	seenActions := map[string]bool{}
	for _, action := range definition.Actions {
		action.Key = normalizeKey(firstNonEmpty(action.Key, action.Label))
		action.Label = firstNonEmpty(action.Label, action.Key)
		action.HandlerRef = strings.TrimSpace(action.HandlerRef)
		if action.Key == "" || seenActions[action.Key] {
			continue
		}
		seenActions[action.Key] = true
		actions = append(actions, action)
	}
	definition.Actions = actions
	return definition
}

func Validate(definition Definition) ValidationErrors {
	definition = Normalize(definition)
	errs := ValidationErrors{}
	if definition.Key == "" {
		errs["key"] = "Flow key is required."
	}
	if len(definition.Steps) == 0 {
		errs["steps"] = "At least one flow step is required."
	}
	for index, action := range definition.Actions {
		if action.HandlerRef == "" {
			errs["actions."+action.Key] = "Action handler ref is required."
			if action.Key == "" {
				errs[fmt.Sprintf("actions.%d", index)] = "Action handler ref is required."
			}
		}
	}
	return errs
}

func Catalog(definitions ...Definition) []Definition {
	out := make([]Definition, 0, len(definitions))
	seen := map[string]bool{}
	for _, definition := range definitions {
		normalized := Normalize(definition)
		if normalized.Key == "" || seen[normalized.Key] {
			continue
		}
		seen[normalized.Key] = true
		out = append(out, normalized)
	}
	return out
}

func Find(catalog []Definition, key string) (Definition, bool) {
	key = normalizeKey(key)
	for _, definition := range catalog {
		if normalizeKey(definition.Key) == key {
			return Normalize(definition), true
		}
	}
	return Definition{}, false
}

func Contact(handlerRef string) Definition {
	return Definition{
		Key:         "contact",
		Label:       "Contact",
		Description: "General contact form.",
		Steps:       []Step{{Key: "message", Label: "Message"}},
		Actions: []Action{{
			Key:        "submit",
			Label:      "Send message",
			HandlerRef: handlerRef,
			Fields: []workbench.Field{
				{Name: "name", Label: "Name", Kind: workbench.FieldText, Required: true},
				{Name: "email", Label: "Email", Kind: workbench.FieldText, Required: true},
				{Name: "message", Label: "Message", Kind: workbench.FieldTextarea, Required: true},
			},
		}},
	}
}

func ScheduleRequest(handlerRef string) Definition {
	return Definition{
		Key:         "schedule-request",
		Label:       "Schedule request",
		Description: "Request a visit, appointment, tour, or class time.",
		Steps:       []Step{{Key: "request", Label: "Request"}},
		Actions: []Action{{
			Key:        "submit",
			Label:      "Request time",
			HandlerRef: handlerRef,
			Fields: []workbench.Field{
				{Name: "guardianName", Label: "Guardian name", Kind: workbench.FieldText, Required: true},
				{Name: "email", Label: "Email", Kind: workbench.FieldText, Required: true},
				{Name: "preferredTime", Label: "Preferred time", Kind: workbench.FieldDateTime},
				{Name: "notes", Label: "Notes", Kind: workbench.FieldTextarea},
			},
		}},
	}
}

func Enrollment(handlerRef string) Definition {
	return Definition{
		Key:         "enrollment",
		Label:       "Enrollment",
		Description: "Program enrollment or waitlist request.",
		Steps:       []Step{{Key: "family", Label: "Family"}, {Key: "child", Label: "Child"}},
		Actions: []Action{{
			Key:        "submit",
			Label:      "Submit enrollment",
			HandlerRef: handlerRef,
			Fields: []workbench.Field{
				{Name: "guardianName", Label: "Guardian name", Kind: workbench.FieldText, Required: true},
				{Name: "childName", Label: "Child name", Kind: workbench.FieldText, Required: true},
				{Name: "childAge", Label: "Child age", Kind: workbench.FieldText, Required: true},
				{Name: "program", Label: "Program", Kind: workbench.FieldRelation},
			},
		}},
	}
}

func normalizeKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
