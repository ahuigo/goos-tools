package nets

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ahuigo/goos-tools/helper"
)

// readNetStat 从对应的文件中读取网络统计值
func readNetStat(interfaceName string, stat string) (int64, error) {
	switch stat {
	case "rx_bytes", "tx_bytes":
		if infBytes, err := getInterfaceTraffic(interfaceName); err != nil {
			return -1, err
		} else {
			if stat == "rx_bytes" {
				return infBytes.Ibytes,nil
			} else {
				return infBytes.Obytes,nil
			}
		}
	}
	return -1, nil
}


func getInterfaceTraffic(interfaceName string) (*InterfaceTraffic, error) {
	out, err := exec.Command("netstat", "-i", "-I", interfaceName, "-b").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected output: %s", out)
	}
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) == 11 {
			return &InterfaceTraffic{
				Name:    fields[0],
				Mtu:     fields[1],
				Network: fields[2],
				Address: fields[3],
				Ipkts:   helper.ParseInt64(fields[4]),
				Ierrs:   helper.ParseInt64(fields[5]),
				Ibytes:  helper.ParseInt64(fields[6]),
				Opkts:   helper.ParseInt64(fields[7]),
				Oerrs:   helper.ParseInt64(fields[8]),
				Obytes:  helper.ParseInt64(fields[9]),
				Coll:    helper.ParseInt64(fields[10]),
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid output: %s", out)
}

type InterfaceTraffic struct {
	Name    string // The name of the network interface
	Mtu     string // The Maximum Transmission Unit (MTU) size
	Network string // The network type
	Address string // The physical (MAC) address of the interface
	Ibytes  int64  // The total number of bytes received on the interface
	Ipkts   int64  // The number of packets received on the interface
	Ierrs   int64  // The number of input errors on the interface
	Obytes  int64  // The total number of bytes **sent** on the interface
	Opkts   int64  // The number of packets sent on the interface
	Oerrs   int64  // The number of output errors on the interface
	Coll    int64  // 发生的碰撞数量。在以太网中，如果两个设备在同一时间发送数据，就会发生碰撞。这里是0，表示没有碰撞。
}

// readSysctl("net.ipv4.tcp_rmem")从 /proc/sys 目录中读取指定的内核参数
// macosx$ sysctl -n kern.ipc.maxsockbuf// 8388608
// macosx$ sysctl -a
var (
	sysctlKeyMap = map[string]string{
		// net.inet.tcp.recvspace定义了TCP接收缓冲区的默认大小，而kern.ipc.maxsockbuf定义了所有类型的socket缓冲区的最大值。
		"net.ipv4.tcp_rmem": "net.inet.tcp.recvspace",
		"net.ipv4.tcp_wmem": "net.inet.tcp.sendspace", //"kern.ipc.maxsockbuf",
		// "net.ipv4.tcp_mem":  "vm.stats.vm.v_page_count",// 所有tcp 堆栈的总内存限制
	}
)

func readSysctl(param string) (string, error) {
	key := sysctlKeyMap[param]
	if key == "" {
		switch param {
		case "net.ipv4.tcp_mem", "net.ipv4.udp_mem", "net.core.rmem_default", "net.core.wmem_default", "net.core.rmem_max", "net.core.wmem_max":

			return "", nil
		default:
			key = param
		}
	}
	return readSysctlRaw(key)
}

func readSysctlRaw(param string) (string, error) {
	cmd := exec.Command("sysctl", "-n", param)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
