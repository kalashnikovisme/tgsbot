package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStandup(t *testing.T) {
	testCases := []struct {
		text      string
		isStandup bool
	}{
		{"hello", false},
		{"yesterday, today, blockers", true},
	}

	for _, tc := range testCases {
		result := isStandup(tc.text)
		assert.Equal(t, tc.isStandup, result)
	}
}
