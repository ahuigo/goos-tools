package nets

import (
	"testing"

	"github.com/ahuigo/goos-tools/helper"
	"github.com/stretchr/testify/assert"
)

func TestGetStats(t *testing.T) {
	// Test case 1: Successful execution
	stats, err := GetStats("")
	assert.NoError(t, err)
	assert.NotEqual(t, "", stats.InterfaceName)
	helper.PrintPretty(stats)
}