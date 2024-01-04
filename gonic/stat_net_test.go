package gonic

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NetStat: show net stat information like shell `netstat`
func TestStatNetTpl(t *testing.T) {
	req :=httptest.NewRequest("GET", "/netstat", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, ctx := mockContext(req)
	NetStat(ctx)
	assert.Equal(t, []string(nil),ctx.Errors.Errors())
	assert.Equal(t, 0, len(ctx.Errors))
	t.Log(resp.Body.String())
}
