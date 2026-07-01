package contentful

import "fmt"

// Region identifies a Contentful infrastructure region.
type Region string

const (
	RegionUS Region = "us"
	RegionEU Region = "eu"
)

// ParseRegion validates s against the known regions and returns the typed Region.
// An empty string is accepted and represents the default (US) region.
func ParseRegion(s string) (Region, error) {
	r := Region(s)
	switch r {
	case "", RegionUS, RegionEU:
		return r, nil
	default:
		return "", fmt.Errorf("unknown region %q: supported values are %q and %q", s, RegionUS, RegionEU)
	}
}
