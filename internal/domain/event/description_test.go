package event

import "testing"

func TestDescription(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		description string
	}{
		{"success new description", true, "description"},
		{"failure empty description", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			description, err := NewDescription(tt.description)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && description.String() != tt.description {
				t.Errorf("String() = %v, want %v", description.String(), tt.description)
			}
		})
	}
}
