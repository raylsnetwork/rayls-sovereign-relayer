package repository_test

import (
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQueryOptions(t *testing.T) {
	tests := []struct {
		name      string
		opts      []repository.Option
		wantLimit int
	}{
		{
			name:      "no options returns default zero limit",
			opts:      nil,
			wantLimit: 0,
		},
		{
			name:      "single WithLimit sets limit",
			opts:      []repository.Option{repository.WithLimit(5)},
			wantLimit: 5,
		},
		{
			name:      "last WithLimit wins",
			opts:      []repository.Option{repository.WithLimit(3), repository.WithLimit(10)},
			wantLimit: 10,
		},
		{
			name:      "boundary value limit of 1",
			opts:      []repository.Option{repository.WithLimit(1)},
			wantLimit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repository.GetQueryOptions(tt.opts...)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantLimit, got.Limit)
		})
	}
}

func TestNewQueryOptions(t *testing.T) {
	got := repository.NewQueryOptions()
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Limit)
}
