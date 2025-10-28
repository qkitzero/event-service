package event

import (
	"fmt"
	"regexp"
)

type Color string

var colorRegexp = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func (c Color) String() string {
	return string(c)
}

func NewColor(s string) (Color, error) {
	if s == "" {
		return Color("#FFFFFF"), nil
	}

	if !colorRegexp.MatchString(s) {
		return Color(""), fmt.Errorf("invalid color")
	}

	return Color(s), nil
}
