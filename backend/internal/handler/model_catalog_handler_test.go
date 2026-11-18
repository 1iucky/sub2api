package handler

import (
	"reflect"
	"testing"
)

func TestParseOptionalInt64s(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []int64
	}{
		{name: "empty", raw: "", want: nil},
		{name: "multiple values", raw: "42, 7,42", want: []int64{42, 7}},
		{name: "ignores invalid values", raw: "42,zero,-3,7", want: []int64{42, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOptionalInt64s(tt.raw); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseOptionalInt64s(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
