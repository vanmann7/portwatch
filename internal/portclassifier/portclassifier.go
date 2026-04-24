// Package portclassifier assigns a category to a port based on its number,
// known service name, and optional user-defined overrides.
package portclassifier

import "fmt"

// Category represents the classification of a port.
type Category string

const (
	CategorySystem    Category = "system"    // 0–1023
	CategoryUser      Category = "user"      // 1024–49151
	CategoryDynamic   Category = "dynamic"   // 49152–65535
	CategoryOverride  Category = "override"  // explicitly set by caller
	CategoryUnknown   Category = "unknown"
)

// Classifier assigns a Category to a port number.
type Classifier struct {
	overrides map[uint16]Category
}

// New returns a Classifier with optional override rules.
// Each override maps a port number to a Category string.
func New(overrides map[uint16]string) (*Classifier, error) {
	c := &Classifier{overrides: make(map[uint16]Category)}
	for port, raw := range overrides {
		cat := Category(raw)
		if !validCategory(cat) {
			return nil, fmt.Errorf("portclassifier: unknown category %q for port %d", raw, port)
		}
		c.overrides[port] = cat
	}
	return c, nil
}

// Classify returns the Category for the given port.
func (c *Classifier) Classify(port uint16) Category {
	if cat, ok := c.overrides[port]; ok {
		return cat
	}
	switch {
	case port < 1024:
		return CategorySystem
	case port < 49152:
		return CategoryUser
	default:
		return CategoryDynamic
	}
}

// ClassifyAll returns a map of port → Category for every port in the slice.
func (c *Classifier) ClassifyAll(ports []uint16) map[uint16]Category {
	out := make(map[uint16]Category, len(ports))
	for _, p := range ports {
		out[p] = c.Classify(p)
	}
	return out
}

func validCategory(cat Category) bool {
	switch cat {
	case CategorySystem, CategoryUser, CategoryDynamic, CategoryOverride, CategoryUnknown:
		return true
	}
	return false
}
