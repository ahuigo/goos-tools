package examples

import (
	"testing"

	"github.com/ahuigo/goos-tools/lsof"
	"github.com/stretchr/testify/assert"
)

func TestGetLsofTcps(t *testing.T) {
	tcps, err := lsof.GetLsofTcps()
	assert.NoError(t, err)
	for _, tcp := range tcps {
		t.Logf("%s, %s, %s, %s, %s\n", tcp.Pid, tcp.Command, tcp.LocalAddr, tcp.ForeignAddr, tcp.TcpState)
	}
}
