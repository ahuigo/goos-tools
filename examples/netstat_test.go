package examples

import (
	"testing"

	"github.com/ahuigo/goos-tools/netstat"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTcpConnections(t *testing.T) {
	tcps, err := netstat.GetAllTcpConnections()
	assert.NoError(t, err)
	for _, tcp := range tcps {
		t.Logf("local:%s, remote:%s, state:%s", tcp.LocalAddr, tcp.ForeignAddr, tcp.State)
	}
}
