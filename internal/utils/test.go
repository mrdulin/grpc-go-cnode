package utils

import "testing"

func MarkedAsIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip(t.Name())
	}
}
