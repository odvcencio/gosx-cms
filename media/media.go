package media

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

type Item struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type Asset struct {
	ID          string     `json:"id"`
	URL         string     `json:"url"`
	Alt         string     `json:"alt"`
	Filename    string     `json:"filename"`
	ContentType string     `json:"contentType"`
	Size        int64      `json:"size"`
	Variants    Variants   `json:"variants,omitempty"`
	FocalX      float64    `json:"focalX,omitempty"`
	FocalY      float64    `json:"focalY,omitempty"`
	ArchivedAt  *time.Time `json:"archivedAt,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

type Variants map[string]Variant

type Variant struct {
	URL         string `json:"url"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

type FocalPoint struct {
	X float64
	Y float64
}

type Usage struct {
	Kind  string
	Title string
	ID    string
	Href  string
}

type Input struct {
	URL         string
	Alt         string
	Filename    string
	ContentType string
	Size        int64
	Variants    Variants
	FocalX      float64
	FocalY      float64
}

type Filter struct {
	Archived *bool
}

type Store interface {
	ListMedia(...Filter) []Asset
	MediaByID(string) (Asset, bool)
	MediaUsage(string) []Usage
	CreateMediaAsset(Input) (Asset, error)
	UpdateMediaAsset(string, Input) (Asset, error)
}

type PickerField struct {
	Name       string
	Label      string
	URL        string
	Alt        string
	Required   bool
	Assets     []Asset
	AssetCount int
	HasAssets  bool
}

func Picker(name, label string, value Item, assets []Asset, required bool) PickerField {
	assets = CloneAssets(assets)
	return PickerField{
		Name:       strings.TrimSpace(name),
		Label:      strings.TrimSpace(label),
		URL:        strings.TrimSpace(value.URL),
		Alt:        strings.TrimSpace(value.Alt),
		Required:   required,
		Assets:     assets,
		AssetCount: len(assets),
		HasAssets:  len(assets) > 0,
	}
}

func NormalizeItem(item Item, fallbackTitle string) Item {
	return Item{
		URL: strings.TrimSpace(item.URL),
		Alt: firstNonEmpty(item.Alt, fallbackTitle),
	}
}

func NormalizeList(items []Item, fallbackTitle string, placeholder bool) []Item {
	out := make([]Item, 0, len(items))
	seen := map[string]bool{}
	for _, item := range items {
		normalized := NormalizeItem(item, fallbackTitle)
		if strings.TrimSpace(normalized.URL) == "" {
			continue
		}
		key := strings.ToLower(normalized.URL)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, normalized)
	}
	if len(out) == 0 && placeholder {
		return []Item{{URL: "/media/placeholder.svg", Alt: firstNonEmpty(fallbackTitle, "Image")}}
	}
	return out
}

func NormalizeAsset(input Input, asset Asset) Asset {
	asset.URL = strings.TrimSpace(input.URL)
	asset.Alt = strings.TrimSpace(input.Alt)
	asset.Filename = strings.TrimSpace(input.Filename)
	asset.FocalX = input.FocalX
	asset.FocalY = input.FocalY
	if strings.TrimSpace(input.ContentType) != "" {
		asset.ContentType = strings.TrimSpace(input.ContentType)
	}
	if input.Size > 0 {
		asset.Size = input.Size
	}
	if input.Variants != nil {
		asset.Variants = NormalizeVariants(input.Variants)
	}
	if asset.Filename == "" {
		asset.Filename = filepath.Base(strings.TrimSpace(input.URL))
	}
	return asset
}

func NormalizeVariants(variants Variants) Variants {
	if len(variants) == 0 {
		return nil
	}
	out := Variants{}
	for name, variant := range variants {
		name = strings.TrimSpace(name)
		variant.URL = strings.TrimSpace(variant.URL)
		variant.ContentType = strings.TrimSpace(variant.ContentType)
		if name == "" || variant.URL == "" {
			continue
		}
		out[name] = variant
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func NormalizeFocalPoint(x, y float64) FocalPoint {
	if x < 0 || x > 1 {
		x = 0
	}
	if y < 0 || y > 1 {
		y = 0
	}
	return FocalPoint{X: x, Y: y}
}

func ParseLines(value string) []Item {
	out := []Item{}
	for _, line := range strings.Split(value, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "|", 2)
		if strings.TrimSpace(parts[0]) == "" {
			continue
		}
		alt := ""
		if len(parts) > 1 {
			alt = strings.TrimSpace(parts[1])
		}
		out = append(out, Item{URL: strings.TrimSpace(parts[0]), Alt: alt})
	}
	return out
}

func FormatLines(items []Item) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		item = NormalizeItem(item, "")
		if item.URL == "" {
			continue
		}
		line := item.URL
		if item.Alt != "" {
			line += " | " + item.Alt
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func CloneAsset(asset Asset) Asset {
	raw, _ := json.Marshal(asset)
	var out Asset
	_ = json.Unmarshal(raw, &out)
	return out
}

func CloneAssets(assets []Asset) []Asset {
	out := make([]Asset, len(assets))
	for i, asset := range assets {
		out[i] = CloneAsset(asset)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
