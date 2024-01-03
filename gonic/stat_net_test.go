package gonic

import (
	"testing"

	"github.com/ahuigo/goos-tools/netstat"
)

// NetStat: show net stat information like shell `netstat`
func TestStatNetTpl(t *testing.T) {
	conns, _ := netstat.GetAllTcpConnections()
	s := formatNetwork(conns)
	t.Log(string(s))

}
