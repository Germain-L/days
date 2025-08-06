package valueobjects

import (
	"errors"
	"regexp"
)

// HexColor represents a valid hex color code
type HexColor struct {
	value string
}

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// NewHexColor creates a new HexColor value object
func NewHexColor(color string) (*HexColor, error) {
	if color == "" {
		return nil, errors.New("hex color cannot be empty")
	}

	if !hexColorRegex.MatchString(color) {
		return nil, errors.New("invalid hex color format, must be #RRGGBB")
	}

	return &HexColor{value: color}, nil
}

// Value returns the string value of the hex color
func (hc *HexColor) Value() string {
	return hc.value
}

// String implements the Stringer interface
func (hc *HexColor) String() string {
	return hc.value
}

// Equals checks if two hex colors are equal
func (hc *HexColor) Equals(other *HexColor) bool {
	if other == nil {
		return false
	}
	return hc.value == other.value
}
