package netstat

import (
	"strings"

	. "github.com/ahuigo/goos-tools"
	"github.com/ahuigo/goos-tools/lsof"
)

func GetAllTcpConnections() (conns []TcpConnection, err error) {
	lsofTcps, err := lsof.GetLsofTcps()
	if err != nil {
		return conns, err
	}
	conns = make([]TcpConnection, len(lsofTcps))
	for i, tcp := range lsofTcps {
		// convert TCP/UDP+IPv6 to tcp6/upd6
		proto:= strings.ToLower(tcp.Node) + IfThen(tcp.Type == "IPv6", "6", "")
		conns[i] = TcpConnection{
			Proto:       proto,
			LocalAddr:   tcp.LocalAddr,
			ForeignAddr: tcp.ForeignAddr,
			State:       tcp.TcpState,
			Pid:         tcp.Pid,
			Program:     tcp.Command,
		}
	}
	return conns, err
}
