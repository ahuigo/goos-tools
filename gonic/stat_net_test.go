package gonic

import (
	"net/http/httptest"
	"testing"
)

// NetStat: show net stat information like shell `netstat`
func TestStatNetTpl(t *testing.T) {
	req :=httptest.NewRequest("GET", "/netstat", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, ctx := mockContext(req)
	NetStat(ctx)
	t.Log(resp.Body.String())
}
