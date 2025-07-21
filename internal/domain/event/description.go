package event

import (
	"fmt"
	"strings"
)

type Description string

func (d Description) String() string {
	return string(d)
}

func NewDescription(s string) (Description, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Description(""), fmt.Errorf("invalid description")
	}
	return Description(s), nil
}
