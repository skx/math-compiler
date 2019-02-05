package token

import (
	"testing"
)

// Test looking up values succeeds, then fails
func TestLookup(t *testing.T) {

	for key, val := range keywords {

		// Obviously this will pass.
		if LookupIdentifier(string(key)) != val {
			t.Errorf("Lookup of %s failed", key)
		}

	}
}
