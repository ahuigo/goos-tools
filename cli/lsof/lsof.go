package lsof

import (
	"errors"
	"strings"

	"github.com/ahuigo/goos-tools/shell"
)

type LsofItem struct {
	Command string
	Pid     string
	User    string
	Fd      string
	Device  string
	//SIZE/OFF 代表常规文件的偏移，对于网络socket无意义。
	SizeOff string
	// IPv6/IPv4/Unix/...
	Type string
	// 文件inode 或 TCP/UDP
	Node string
	// 文件或socket链接
	Name        string
	LocalAddr   string
	ForeignAddr string
	// LISTEN/SYN_SENT/ESTABLISHED/CLOSE_WAIT/CLOSING/LAST_ACK/TIME_WAIT/CLOSED
	TcpState string
}

/*
*
lsof输出保存到stdout中，格式如下：
```
COMMAND     PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
postgres    530 ahui    7u  IPv6 0x803b406e0f0bace3      0t0  TCP [::1]:5432 (LISTEN)
postgres    530 ahui    8u  IPv4 0x803b407c764b7f2b      0t0  TCP 127.0.0.1:5432 (LISTEN)
```
*/
func GetLsofTcps() (out []LsofItem, err error) {
	// lsof -iTCP -sTCP:ESTABLISHED,LISTEN -n -P
	cmd := "lsof -iTCP -n -P"
	stdout, errmsg, errno := shell.ExecCommand("sh", "-c", cmd)
	if errno != 0 {
		err = errors.New(errmsg)
		return out, err
	}
	lines := strings.Split(stdout, "\n")
	if len(lines) <= 1 {
		return out, nil
	}
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 10 {
			continue
		}
		addr := fields[8]
		state := strings.Trim(fields[9], "()") // "(LISTEN)"
		localAddr, foreignAddr, _ := strings.Cut(addr, "->")

		out = append(out, LsofItem{
			Command:     fields[0],
			Pid:         fields[1],
			User:        fields[2],
			Fd:          fields[3],
			Type:        fields[4],
			Device:      fields[5],
			SizeOff:     fields[6],
			Node:        fields[7],
			Name:        fields[8],
			LocalAddr:   localAddr,
			ForeignAddr: foreignAddr,
			TcpState:    state, // LISTEN/SYN_SENT/ESTABLISHED/CLOSE_WAIT/CLOSING/LAST_ACK/TIME_WAIT/CLOSED
		})

	}
	return out, nil
}
