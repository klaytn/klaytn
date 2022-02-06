package storage

import (
	"os"
	"testing"
)

func SkipLocalTest(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping testing in Local environment")
	}
}
