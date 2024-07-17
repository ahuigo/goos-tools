# goos-tools
[![tag](https://img.shields.io/github/tag/ahuigo/goos-tools.svg)](https://github.com/ahuigo/goos-tools/tags)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/goos-tools?status.svg)](https://pkg.go.dev/github.com/ahuigo/goos-tools)
![Build Status](https://github.com/ahuigo/goos-tools/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/ahuigo/goos-tools)](https://goreportcard.com/report/github.com/ahuigo/goos-tools)
[![Coverage](https://img.shields.io/codecov/c/github/ahuigo/goos-tools)](https://codecov.io/gh/ahuigo/goos-tools)
[![Contributors](https://img.shields.io/github/contributors/ahuigo/goos-tools)](https://github.com/ahuigo/goos-tools/graphs/contributors)
[![License](https://img.shields.io/github/license/ahuigo/goos-tools)](./LICENSE)


Os tools for golang

- [goos-tools](#goos-tools)
  - [Net statistic](#net-statistic)
  - [lsof](#lsof)
- [gonic-tools](#gonic-tools)
  - [Show Stat](#show-stat)

## Net statistic
Example1: https://github.com/ahuigo/goos-tools/blob/main/nets/stats_linux_test.go
    import (
        "testing"
        "github.com/ahuigo/goos-tools/helper"
        "github.com/stretchr/testify/assert"
    )

    func TestGetStats(t *testing.T) {
        stats, err := GetStats("")
        assert.NoError(t, err)
        assert.NotEqual(t, "", stats.InterfaceName)
        helper.PrintPretty(stats)
    }

Example2: https://github.com/ahuigo/goos-tools/blob/master/examples

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
