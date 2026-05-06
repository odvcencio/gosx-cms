package blocks

import (
	"testing"

	"github.com/odvcencio/gosx-admin/blockstudio"
)

func TestHomeCatalogNormalizes(t *testing.T) {
	blocks := blockstudio.Normalize(nil, HomeCatalog())
	if len(blocks) != 4 {
		t.Fatalf("expected four home blocks, got %#v", blocks)
	}
	for _, block := range blocks {
		if !block.Enabled {
			t.Fatalf("default home block should be enabled: %#v", block)
		}
	}
	if !blockstudio.KeyAllowed("hero", HomeCatalog()) {
		t.Fatal("expected hero block to be allowed")
	}
}
