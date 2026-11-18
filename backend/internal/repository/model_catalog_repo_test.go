package repository

import (
	"reflect"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestBuildModelCatalogWhereFiltersMultipleVendors(t *testing.T) {
	where, args := buildModelCatalogWhere(service.ModelCatalogListFilters{
		VendorIDs: []int64{42, 7},
	})

	if !reflect.DeepEqual(args, []any{int64(42), int64(7)}) {
		t.Fatalf("args = %v, want [42 7]", args)
	}
	if got, want := where[len(where)-1], "mc.vendor_id IN ($1, $2)"; got != want {
		t.Fatalf("vendor filter = %q, want %q", got, want)
	}
}
