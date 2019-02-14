package token

import (
	"testing"
)

// Test looking up values succeeds.
func TestLookup(t *testing.T) {

	for key, val := range keywords {

		// Obviously this will pass.
		if LookupIdentifier(string(key)) != val {
			t.Errorf("Lookup of %s failed", key)
		}

	}
}

// Test looking up unknown-values.
func TestLookupFailures(t *testing.T) {

	keywords := []string{"foo", "bar", "baz"}

	for _, val := range keywords {

		// Obviously this will pass.
		if LookupIdentifier(string(val)) != ERROR {
			t.Errorf("Lookup of %s was expected to fail, but didn't", val)
		}

	}
}
