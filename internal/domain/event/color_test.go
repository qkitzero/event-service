package event

import "testing"

func TestNewColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		color         string
		expectedColor string
	}{
		{"success new color", true, "#FF0000", "#FF0000"},
		{"success default color", true, "", "#FFFFFF"},
		{"failure invalid color", false, "red", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			color, err := NewColor(tt.color)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && color.String() != tt.expectedColor {
				t.Errorf("String() = %v, want %v", color.String(), tt.expectedColor)
			}
		})
	}
}
