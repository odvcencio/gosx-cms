package studio

import (
	"fmt"
	gosxstudio "m31labs.dev/gosx-studio"
	"strings"
)

type Assessment struct {
	Status       gosxstudio.ShellReadinessStatus
	Summary      string
	Detail       string
	Value        string
	Count        int
	Total        int
	Missing      []string
	MissingCount int
}

type RequiredField struct {
	Label string
	Value string
}

type RequiredFieldOptions struct {
	ReadySummary        string
	ReadyDetail         string
	MissingStatus       gosxstudio.ShellReadinessStatus
	MissingSummary      string
	MissingDetailPrefix string
}

type FlowExecutionOptions struct {
	ReadyDetail   string
	PartialDetail string
	EmptyDetail   string
}

func AssessRequiredFields(fields []RequiredField, options RequiredFieldOptions) Assessment {
	missing := MissingRequiredFields(fields)
	status := gosxstudio.ShellReadinessReady
	summary := gosxstudio.FirstNonEmpty(options.ReadySummary, "Required fields present")
	detail := gosxstudio.FirstNonEmpty(options.ReadyDetail, "Required fields are filled.")
	if len(missing) > 0 {
		status = gosxstudio.NormalizeShellReadinessStatus(options.MissingStatus)
		summary = gosxstudio.FirstNonEmpty(options.MissingSummary, fmt.Sprintf("%d missing fields", len(missing)))
		prefix := gosxstudio.FirstNonEmpty(options.MissingDetailPrefix, "Missing")
		detail = strings.TrimSpace(prefix) + " " + strings.Join(missing, ", ") + "."
	}
	return Assessment{
		Status:       status,
		Summary:      summary,
		Detail:       detail,
		Value:        summary,
		Count:        len(fields) - len(missing),
		Total:        len(fields),
		Missing:      missing,
		MissingCount: len(missing),
	}
}

func MissingRequiredFields(fields []RequiredField) []string {
	missing := []string{}
	for _, field := range fields {
		label := strings.TrimSpace(field.Label)
		if label == "" {
			continue
		}
		if strings.TrimSpace(field.Value) == "" {
			missing = append(missing, label)
		}
	}
	return missing
}

func AssessExecutableFlows(total, executable int, options FlowExecutionOptions) Assessment {
	if total < 0 {
		total = 0
	}
	if executable < 0 {
		executable = 0
	}
	if executable > total {
		executable = total
	}
	status := gosxstudio.ShellReadinessReady
	detail := gosxstudio.FirstNonEmpty(options.ReadyDetail, "Registered flows have executable handler refs.")
	if total == 0 || executable == 0 {
		status = gosxstudio.ShellReadinessNext
		detail = gosxstudio.FirstNonEmpty(options.EmptyDetail, "Register executable flows before publishing public forms.")
	} else if executable < total {
		status = gosxstudio.ShellReadinessWatch
		detail = gosxstudio.FirstNonEmpty(options.PartialDetail, "Some registered flows still need handler refs.")
	}
	summary := fmt.Sprintf("%d/%d executable", executable, total)
	return Assessment{
		Status:  status,
		Summary: summary,
		Detail:  detail,
		Value:   summary,
		Count:   executable,
		Total:   total,
	}
}

func ExecutableFlowCardCount(flows []FlowCard) int {
	count := 0
	for _, flow := range normalizeFlowCards(flows) {
		if flow.CanExecute {
			count++
		}
	}
	return count
}
