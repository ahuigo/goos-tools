# goos-utils
Os utils for golang

## netstat
For more examples, refer to [examples](https://github.com/ahuigo/goos-utils/blob/master/examples) 

    import (
        "testing"
        "github.com/ahuigo/goos-utils/netstat"
        "github.com/stretchr/testify/assert"
    )

    func TestGetAllTcpConnections(t *testing.T) {
        tcps, err := netstat.GetAllTcpConnections()
        assert.NoError(t, err)
        for _, tcp := range tcps {
            t.Logf("local:%s, remote:%s, state:%s", tcp.LocalAddr, tcp.ForeignAddr, tcp.State)
        }
    }

## lsof
For more examples, refer to [examples](https://github.com/ahuigo/goos-utils/blob/master/examples) 

```
import (
	"testing"

	"github.com/ahuigo/goos-utils/lsof"
	"github.com/stretchr/testify/assert"
)

func TestGetLsofTcps(t *testing.T) {
	tcps, err := lsof.GetLsofTcps()
	assert.NoError(t, err)
	for _, tcp := range tcps {
		t.Logf("%s, %s, %s, %s, %s\n", tcp.Pid, tcp.Command, tcp.LocalAddr, tcp.ForeignAddr, tcp.TcpState)
	}
}

```
