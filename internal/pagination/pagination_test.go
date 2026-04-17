package pagination_test

import (
	"testing"

	"github.com/fvdm-otinga/fireflies-cli/internal/pagination"
)

func TestNextPage(t *testing.T) {
	tests := []struct {
		name          string
		skip          int
		limit         int
		returnedCount int
		wantNext      *int
	}{
		{
			name:          "full page — more results available",
			skip:          0,
			limit:         10,
			returnedCount: 10,
			wantNext:      intPtr(10),
		},
		{
			name:          "partial page — last page",
			skip:          0,
			limit:         10,
			returnedCount: 7,
			wantNext:      nil,
		},
		{
			name:          "empty page",
			skip:          0,
			limit:         10,
			returnedCount: 0,
			wantNext:      nil,
		},
		{
			name:          "second full page",
			skip:          10,
			limit:         10,
			returnedCount: 10,
			wantNext:      intPtr(20),
		},
		{
			name:          "zero limit means no next",
			skip:          0,
			limit:         0,
			returnedCount: 50,
			wantNext:      nil,
		},
		{
			name:          "limit 1 full page",
			skip:          5,
			limit:         1,
			returnedCount: 1,
			wantNext:      intPtr(6),
		},
		{
			name:          "limit 50 transcript cap — full",
			skip:          0,
			limit:         50,
			returnedCount: 50,
			wantNext:      intPtr(50),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := pagination.NextPage(tc.skip, tc.limit, tc.returnedCount)
			if tc.wantNext == nil && got != nil {
				t.Errorf("expected nil, got %d", *got)
			} else if tc.wantNext != nil && got == nil {
				t.Errorf("expected %d, got nil", *tc.wantNext)
			} else if tc.wantNext != nil && got != nil && *got != *tc.wantNext {
				t.Errorf("expected %d, got %d", *tc.wantNext, *got)
			}
		})
	}
}

func TestNewCursor(t *testing.T) {
	c := pagination.NewCursor(0, 10, 10)
	if c.Skip != 0 || c.Limit != 10 {
		t.Errorf("unexpected Skip/Limit: %d/%d", c.Skip, c.Limit)
	}
	if c.NextSkip == nil || *c.NextSkip != 10 {
		t.Error("expected NextSkip = 10")
	}

	c2 := pagination.NewCursor(10, 10, 5)
	if c2.NextSkip != nil {
		t.Errorf("expected nil NextSkip, got %d", *c2.NextSkip)
	}
}

func intPtr(n int) *int { return &n }
