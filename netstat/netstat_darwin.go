package netstat

import "github.com/ahuigo/goos-tools/lsof"

func GetAllTcpConnections() (conns []TcpConnection, err error) {
	lsofTcps, err := lsof.GetLsofTcps()
	if err != nil {
		return conns, err
	}
	conns = make([]TcpConnection, len(lsofTcps))
	for i, tcp := range lsofTcps {
		conns[i] = TcpConnection{
			Proto:       tcp.Node,
			LocalAddr:   tcp.LocalAddr,
			ForeignAddr: tcp.ForeignAddr,
			State:       tcp.TcpState,
			Pid:         tcp.Pid,
			Program:     tcp.Command,
		}
	}
	return conns, err
}
