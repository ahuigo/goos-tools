# goos-tools
Os tools for golang

## netstat
For more examples, refer to [examples](https://github.com/ahuigo/goos-tools/blob/master/examples) 

    import (
        "testing"
        "github.com/ahuigo/goos-tools/cli/netstat"
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
For more examples, refer to [examples](https://github.com/ahuigo/goos-tools/blob/master/examples) 

```
import (
	"testing"

	"github.com/ahuigo/goos-tools/cli/lsof"
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

# gonic-tools
Tools for gonic

## Show Stat
Add routes:

    import (
        "github.com/ahuigo/goos-tools/gonic"
    )

    r := gin.New()
    r.GET("/stat/os", gonic.OsStat)
    r.GET("/stat/net", gonic.StatNet)

The code above has functions:
- Show network stat(all tcp connects): http://your.domain/stat/net
- Show os stat(memory and cpu usage): http://your.domain/stat/os
