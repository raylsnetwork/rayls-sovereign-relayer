package proofgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDivCeil(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "7/3 rounds up", a: 7, b: 3, want: 3},
		{name: "6/3 exact", a: 6, b: 3, want: 2},
		{name: "1/1", a: 1, b: 1, want: 1},
		{name: "0/3", a: 0, b: 3, want: 0},
		{name: "5/2 rounds up", a: 5, b: 2, want: 3},
		{name: "1/3 rounds up to 1", a: 1, b: 3, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := divCeil(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("panics on division by zero", func(t *testing.T) {
		assert.Panics(t, func() { divCeil(5, 0) })
	})
}
