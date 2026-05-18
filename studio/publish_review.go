package studio

import (
	"fmt"
	"strings"

	"github.com/odvcencio/gosx"
)

type PublishReview struct {
	Key                string
	ResourceKind       string
	ResourceID         string
	Title              string
	Summary            string
	Status             ReadinessStatus
	Checks             []PublishCheck
	Impacts            []PublishImpact
	PrimaryHref        string
	PrimaryActionLabel string
}

type PublishCheck struct {
	Key         string
	Label       string
	Scope       string
	Status      ReadinessStatus
	Summary     string
	Detail      string
	Href        string
	ActionLabel string
}

type PublishImpact struct {
	Key    string
	Label  string
	Scope  string
	Value  string
	Detail string
	Status ReadinessStatus
}

type PublishReviewOptions struct {
	Class       string
	Kicker      string
	Title       string
	EmptyTitle  string
	EmptyDetail string
	ImpactTitle string
}

func NewPublishReview(checks ...PublishCheck) PublishReview {
	return NormalizePublishReview(PublishReview{Checks: checks})
}

func NewPublishCheck(key, label, scope string, status ReadinessStatus, summary, detail string) PublishCheck {
	return PublishCheck{
		Key:     key,
		Label:   label,
		Scope:   scope,
		Status:  status,
		Summary: summary,
		Detail:  detail,
	}
}

func (check PublishCheck) WithHref(href string) PublishCheck {
	check.Href = href
	return check
}

func (check PublishCheck) WithActionLabel(label string) PublishCheck {
	check.ActionLabel = label
	return check
}

func NewPublishImpact(key, label, scope, value, detail string, status ReadinessStatus) PublishImpact {
	return PublishImpact{
		Key:    key,
		Label:  label,
		Scope:  scope,
		Value:  value,
		Detail: detail,
		Status: status,
	}
}

func NormalizePublishReview(review PublishReview) PublishReview {
	review.Key = normalizeKey(review.Key)
	review.ResourceKind = strings.TrimSpace(review.ResourceKind)
	review.ResourceID = strings.TrimSpace(review.ResourceID)
	review.Title = strings.TrimSpace(review.Title)
	review.Summary = strings.TrimSpace(review.Summary)
	review.PrimaryHref = strings.TrimSpace(review.PrimaryHref)
	review.PrimaryActionLabel = strings.TrimSpace(review.PrimaryActionLabel)
	review.Checks = normalizePublishChecks(review.Checks)
	review.Impacts = normalizePublishImpacts(review.Impacts)
	if review.Key == "" {
		review.Key = normalizeKey(firstNonEmpty(review.ResourceID, review.Title, review.ResourceKind, "publish-review"))
	}
	if review.ResourceKind == "" {
		review.ResourceKind = "Site"
	}
	if review.Title == "" {
		review.Title = "Publish review"
	}
	review.Status = normalizeReadinessStatus(review.Status)
	if review.Status == ReadinessWatch && len(review.Checks) > 0 {
		review.Status = derivedPublishStatus(review.Checks)
	}
	if review.Summary == "" {
		ready, _, _, total := publishCheckCounts(review.Checks)
		review.Summary = fmt.Sprintf("%d/%d clear", ready, total)
	}
	if review.PrimaryActionLabel == "" {
		review.PrimaryActionLabel = readinessActionLabel(review.Status)
	}
	return review
}

func (review PublishReview) Counts() (ready, watch, next, total int) {
	return publishCheckCounts(normalizePublishChecks(review.Checks))
}

func publishCheckCounts(checks []PublishCheck) (ready, watch, next, total int) {
	for _, check := range checks {
		switch check.Status {
		case ReadinessReady:
			ready++
		case ReadinessNext:
			next++
		default:
			watch++
		}
	}
	return ready, watch, next, len(checks)
}

func (review PublishReview) CountSummary() string {
	ready, _, _, total := review.Counts()
	return fmt.Sprintf("%d/%d clear", ready, total)
}

func PublishReviewView(review PublishReview) map[string]any {
	review = NormalizePublishReview(review)
	ready, watch, next, total := review.Counts()
	return map[string]any{
		"key":                review.Key,
		"resourceKind":       review.ResourceKind,
		"resourceID":         review.ResourceID,
		"title":              review.Title,
		"summary":            review.Summary,
		"status":             string(review.Status),
		"statusLabel":        readinessStatusLabel(review.Status),
		"readyCount":         ready,
		"watchCount":         watch,
		"nextCount":          next,
		"total":              total,
		"checks":             publishCheckViews(review.Checks),
		"impacts":            publishImpactViews(review.Impacts),
		"primaryHref":        review.PrimaryHref,
		"hasPrimaryHref":     review.PrimaryHref != "",
		"primaryActionLabel": review.PrimaryActionLabel,
	}
}

func RenderPublishReviewPanel(review PublishReview, options PublishReviewOptions) gosx.Node {
	className := firstNonEmpty(options.Class, "studio-publish-review")
	review = NormalizePublishReview(review)
	ready, watch, next, _ := review.Counts()
	children := []gosx.Node{
		gosx.El("div", gosx.Attrs(gosx.Attr("class", className+"__head")),
			gosx.El("div", nil,
				gosx.El("p", gosx.Attrs(gosx.Attr("class", className+"__kicker")), gosx.Text(firstNonEmpty(options.Kicker, review.ResourceKind))),
				gosx.El("h2", nil, gosx.Text(firstNonEmpty(options.Title, review.Title))),
			),
			gosx.El("output", gosx.Attrs(gosx.Attr("class", className+"__count")), gosx.Text(review.CountSummary())),
		),
		gosx.El("p", gosx.Attrs(gosx.Attr("class", className+"__summary-text")), gosx.Text(review.Summary)),
	}
	if len(review.Checks) == 0 {
		children = append(children, gosx.El("article", gosx.Attrs(gosx.Attr("class", className+"__empty")),
			gosx.El("strong", nil, gosx.Text(firstNonEmpty(options.EmptyTitle, "No publish checks"))),
			gosx.El("p", nil, gosx.Text(firstNonEmpty(options.EmptyDetail, "Register content, approval, flow, and deployment checks before publish."))),
		))
	} else {
		children = append(children, gosx.El("div", gosx.Attrs(
			gosx.Attr("class", className+"__summary"),
			gosx.Attr("aria-label", "Publish review summary"),
		),
			gosx.El("span", nil, gosx.El("strong", nil, gosx.Text(fmt.Sprint(ready))), gosx.Text(" clear")),
			gosx.El("span", nil, gosx.El("strong", nil, gosx.Text(fmt.Sprint(watch))), gosx.Text(" watch")),
			gosx.El("span", nil, gosx.El("strong", nil, gosx.Text(fmt.Sprint(next))), gosx.Text(" next")),
		))
		items := make([]gosx.Node, 0, len(review.Checks))
		for _, check := range review.Checks {
			items = append(items, renderPublishCheck(className, check))
		}
		children = append(children, gosx.El("div", gosx.Attrs(gosx.Attr("class", className+"__list")), gosx.Fragment(items...)))
	}
	if len(review.Impacts) > 0 {
		children = append(children, renderPublishImpacts(className, firstNonEmpty(options.ImpactTitle, "Publish impact"), review.Impacts))
	}
	if review.PrimaryHref != "" {
		children = append(children, gosx.El("a", gosx.Attrs(
			gosx.Attr("class", className+"__primary"),
			gosx.Attr("href", review.PrimaryHref),
			gosx.Attr("data-gosx-link", "true"),
		), gosx.Text(review.PrimaryActionLabel)))
	}
	return gosx.El("section", gosx.Attrs(
		gosx.Attr("class", className),
		gosx.Attr("data-studio-publish-review", review.Key),
		gosx.Attr("data-studio-publish-status", string(review.Status)),
	), gosx.Fragment(children...))
}

func normalizePublishChecks(checks []PublishCheck) []PublishCheck {
	out := make([]PublishCheck, 0, len(checks))
	for _, check := range checks {
		check.Key = normalizeKey(check.Key)
		check.Label = strings.TrimSpace(check.Label)
		check.Scope = strings.TrimSpace(check.Scope)
		check.Status = normalizeReadinessStatus(check.Status)
		check.Summary = strings.TrimSpace(check.Summary)
		check.Detail = strings.TrimSpace(check.Detail)
		check.Href = strings.TrimSpace(check.Href)
		check.ActionLabel = strings.TrimSpace(check.ActionLabel)
		if check.Key == "" {
			check.Key = normalizeKey(check.Label)
		}
		if check.Scope == "" {
			check.Scope = "Site"
		}
		if check.Summary == "" {
			check.Summary = readinessStatusLabel(check.Status)
		}
		if check.ActionLabel == "" {
			check.ActionLabel = readinessActionLabel(check.Status)
		}
		if check.Key == "" || check.Label == "" {
			continue
		}
		out = append(out, check)
	}
	return out
}

func normalizePublishImpacts(impacts []PublishImpact) []PublishImpact {
	out := make([]PublishImpact, 0, len(impacts))
	for _, impact := range impacts {
		impact.Key = normalizeKey(impact.Key)
		impact.Label = strings.TrimSpace(impact.Label)
		impact.Scope = strings.TrimSpace(impact.Scope)
		impact.Value = strings.TrimSpace(impact.Value)
		impact.Detail = strings.TrimSpace(impact.Detail)
		impact.Status = normalizeReadinessStatus(impact.Status)
		if impact.Key == "" {
			impact.Key = normalizeKey(impact.Label)
		}
		if impact.Scope == "" {
			impact.Scope = "Site"
		}
		if impact.Value == "" {
			impact.Value = "Tracked"
		}
		if impact.Key == "" || impact.Label == "" {
			continue
		}
		out = append(out, impact)
	}
	return out
}

func derivedPublishStatus(checks []PublishCheck) ReadinessStatus {
	status := ReadinessReady
	for _, check := range checks {
		switch normalizeReadinessStatus(check.Status) {
		case ReadinessNext:
			return ReadinessNext
		case ReadinessWatch:
			status = ReadinessWatch
		}
	}
	return status
}

func publishCheckViews(checks []PublishCheck) []map[string]any {
	out := make([]map[string]any, 0, len(checks))
	for _, check := range checks {
		status := normalizeReadinessStatus(check.Status)
		out = append(out, map[string]any{
			"key":         check.Key,
			"label":       check.Label,
			"scope":       check.Scope,
			"status":      string(status),
			"statusLabel": readinessStatusLabel(status),
			"class":       "studio-publish-review__card studio-publish-review__card--" + string(status),
			"summary":     check.Summary,
			"detail":      check.Detail,
			"href":        check.Href,
			"hasHref":     check.Href != "",
			"actionLabel": firstNonEmpty(check.ActionLabel, readinessActionLabel(status)),
		})
	}
	return out
}

func publishImpactViews(impacts []PublishImpact) []map[string]any {
	out := make([]map[string]any, 0, len(impacts))
	for _, impact := range impacts {
		status := normalizeReadinessStatus(impact.Status)
		out = append(out, map[string]any{
			"key":         impact.Key,
			"label":       impact.Label,
			"scope":       impact.Scope,
			"value":       impact.Value,
			"detail":      impact.Detail,
			"status":      string(status),
			"statusLabel": readinessStatusLabel(status),
			"class":       "studio-publish-review__impact studio-publish-review__impact--" + string(status),
		})
	}
	return out
}

func renderPublishCheck(className string, check PublishCheck) gosx.Node {
	children := []gosx.Node{
		gosx.El("div", gosx.Attrs(gosx.Attr("class", className+"__card-head")),
			gosx.El("div", nil,
				gosx.El("strong", nil, gosx.Text(check.Label)),
				gosx.El("span", nil, gosx.Text(check.Scope)),
			),
			gosx.El("output", nil, gosx.Text(readinessStatusLabel(check.Status))),
		),
		gosx.El("p", gosx.Attrs(gosx.Attr("class", className+"__check-summary")), gosx.Text(check.Summary)),
	}
	if check.Detail != "" {
		children = append(children, gosx.El("p", gosx.Attrs(gosx.Attr("class", className+"__detail")), gosx.Text(check.Detail)))
	}
	if check.Href != "" {
		children = append(children, gosx.El("a", gosx.Attrs(
			gosx.Attr("href", check.Href),
			gosx.Attr("data-gosx-link", "true"),
		), gosx.Text(check.ActionLabel)))
	}
	return gosx.El("article", gosx.Attrs(
		gosx.Attr("class", className+"__card "+className+"__card--"+string(check.Status)),
		gosx.Attr("data-studio-publish-check", check.Key),
	), gosx.Fragment(children...))
}

func renderPublishImpacts(className, title string, impacts []PublishImpact) gosx.Node {
	items := make([]gosx.Node, 0, len(impacts))
	for _, impact := range impacts {
		children := []gosx.Node{
			gosx.El("div", nil,
				gosx.El("strong", nil, gosx.Text(impact.Label)),
				gosx.El("span", nil, gosx.Text(impact.Scope)),
			),
			gosx.El("output", nil, gosx.Text(impact.Value)),
		}
		if impact.Detail != "" {
			children = append(children, gosx.El("p", nil, gosx.Text(impact.Detail)))
		}
		items = append(items, gosx.El("article", gosx.Attrs(
			gosx.Attr("class", className+"__impact "+className+"__impact--"+string(impact.Status)),
			gosx.Attr("data-studio-publish-impact", impact.Key),
		), gosx.Fragment(children...)))
	}
	return gosx.El("div", gosx.Attrs(gosx.Attr("class", className+"__impacts")),
		gosx.El("h3", nil, gosx.Text(title)),
		gosx.El("div", gosx.Attrs(gosx.Attr("class", className+"__impact-list")), gosx.Fragment(items...)),
	)
}
