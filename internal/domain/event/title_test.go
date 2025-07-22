package event

import "testing"

func TestNewTitle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		title   string
	}{
		{"success new title", true, "title"},
		{"failure empty title", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			title, err := NewTitle(tt.title)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && title.String() != tt.title {
				t.Errorf("String() = %v, want %v", title.String(), tt.title)
			}
		})
	}
}
