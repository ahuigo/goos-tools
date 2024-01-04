package netstat

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ahuigo/goos-tools/shell"

	"github.com/samber/lo"
)

/*
*
# netstat -antup 输出保存到stdout中，格式如下：
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 10.244.64.176:54764     192.168.145.102:4318    ESTABLISHED 1/mpush-server
tcp        0      0 10.244.64.176:45146     47.93.126.196:443       ESTABLISHED 1/mpush-server
tcp        0      0 10.244.64.176:43816     10.245.218.86:80        ESTABLISHED 1/mpush-server
*/
func GetAllTcpConnections() (conns []TcpConnection, err error) {
	cmd := "netstat -antp"
	stdout, errmsg, errno := shell.ExecCommand("sh", "-c", cmd)
	if errno != 0 {
		err = errors.New(errmsg)
		return conns, err
	}
	lines := strings.Split(stdout, "\n")
	_, index, ok := lo.FindIndexOf(lines, func(s string) bool {
		return strings.HasPrefix(s, "Proto Recv-Q")
	})
	if len(lines) <= 2 || !ok {
		return conns, nil
	}
	index++

	conns = make([]TcpConnection, 0, len(lines)-index)
	for _, line := range lines[index:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		recvQ, _ := strconv.Atoi(fields[1])
		sendQ, _ := strconv.Atoi(fields[2])
		localAddr := fields[3]
		remoteAddr := fields[4]
		state := fields[5]
		pid, programName, _ := strings.Cut(fields[6], "/")
		if !strings.Contains(fields[6], "/") {
			pid = ""
		}
		conns = append(conns, TcpConnection{
			Proto:       fields[0],
			RecvQ:       recvQ,
			SendQ:       sendQ,
			LocalAddr:   localAddr,
			ForeignAddr: remoteAddr,
			State:       state,
			Pid:         pid,
			Program:     programName,
		})
	}
	return conns, err
}
