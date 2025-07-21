package event

import (
	"fmt"
	"strings"
)

type Title string

func (t Title) String() string {
	return string(t)
}

func NewTitle(s string) (Title, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Title(""), fmt.Errorf("invalid title")
	}
	return Title(s), nil
}
