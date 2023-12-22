package netstats

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

/*
在 /sys/class/net/eth0/statistics/ 目录下，有若干文件记录了网络接口 eth0 的输入输出统计数据。以下是其中一些与网络I/O瓶颈相关的重要文件及其含义：

	rx_bytes 和 tx_bytes: 分别表示接收和发送的总字节数。
	rx_errors 和 tx_errors: 接收和发送时发生错误的数据包数。

分析判断网络I/O是否到达瓶颈可以遵循以下步骤：

	基线测量：获取正常工作负载下的网络统计数据作为基线。
	监控和对比：定期监测当前的统计数据与基线对比，注意任何显著的变化。
	带宽利用率：检查 rx_bytes 和 tx_bytes 以确定实际使用的带宽是否接近网络接口的设计极限。
	错误和丢包：过多的 rx_errors、tx_errors、rx_dropped 或 tx_dropped 可能指示存在问题。
	固件和驱动更新：确保网卡固件和驱动程序是最新的，旧版本可能存在性能问题。
	硬件检查：如果出现较多的 rx_frame_errors、tx_carrier_errors 等物理层错误，可能需要检查网线、交换机端口和其他硬件组件。
	系统资源：确认服务器的CPU、内存等资源是否充裕。瓶颈有时候是由于服务器的处理能力不足造成的。
	配置优化：调整网卡和操作系统的网络设置，比如增大缓冲区大小，改善中断分配（IRQ），开启或关闭TCP卸载特性等。

*/

/*
查看TCP和UDP缓冲区大小

	运行以下命令以查看当前的TCP和UDP缓冲区大小配置：
		# 查看TCP缓冲区大小
		sysctl net.ipv4.tcp_rmem
		sysctl net.ipv4.tcp_wmem

		# 查看UDP缓冲区大小
		sysctl net.core.rmem_default
		sysctl net.core.wmem_default
		sysctl net.core.rmem_max
		sysctl net.core.wmem_max

	这些参数的含义如下：
		tcp_rmem：控制TCP接收缓冲区的最小值、默认值和最大值（字节为单位）。
		tcp_wmem：控制TCP发送缓冲区的最小值、默认值和最大值（字节为单位）。
		rmem_default 和 wmem_default：分别为UDP接收和发送操作设置默认缓冲区大小。
		rmem_max 和 wmem_max：分别设置UDP接收和发送操作的最大缓冲区大小。

增大TCP和UDP缓冲区大小

	要永久更改这些值，请编辑 /etc/sysctl.conf 文件，并添加或修改相应的行。例如：

		# 设置TCP缓冲区大小
		net.ipv4.tcp_rmem = 4096 87380 6291456
		net.ipv4.tcp_wmem = 4096 16384 4194304

		# 设置UDP缓冲区大小
		net.core.rmem_default = 262144
		net.core.wmem_default = 262144
		net.core.rmem_max = 4194304
		net.core.wmem_max = 4194304
		编辑文件后，您需要运行 sysctl -p 来应用更改。

	如果你只想临时更改这些值（在下次重启之前），可以使用 sysctl 命令直接设置它们：

		# 临时设置TCP缓冲区大小
		sudo sysctl -w net.ipv4.tcp_rmem="4096 87380 6291456"
		sudo sysctl -w net.ipv4.tcp_wmem="4096 16384 4194304"

		# 临时设置UDP缓冲区大小
		sudo sysctl -w net.core.rmem_default=262144
		sudo sysctl -w net.core.wmem_default=262144
		sudo sysctl -w net.core.rmem_max=4194304
		sudo sysctl -w net.core.wmem_max=4194304

*/

// NetStats 代表了网络接口的统计数据
type NetStats struct {
	//如果你看到这两个数值在不断增长，并且接近网卡的最大带宽，那可能意味着网络接口接近或达到其容量极限。
	RxBytes int64 `json:"rx_bytes"` // 接收的总字节数
	TxBytes int64 `json:"tx_bytes"` // 发送的总字节数
	// 高错误率可能表明物理连接问题、溢出的缓冲区或配置错误等，这些都可能导致网络性能下降。
	RxErrors int64 `json:"rx_errors"` // 接收时发生错误的数据包数
	TxErrors int64 `json:"tx_errors"` // 发送时发生错误的数据包数
	// rx_dropped 和 tx_dropped: 接收和发送时被丢弃的数据包数。一个高的 rx_dropped 计数可能表示接收缓冲区满了，而系统处理不过来接收的数据包。
	RxDropped int64 `json:"rx_dropped"` // 接收时被丢弃的数据包数
	TxDropped int64 `json:"tx_dropped"` // 发送时被丢弃的数据包数
	// rx_fifo_errors: 表示因为FIFO缓冲区溢出而丢失的数据包数量。如果这个值很高，可能意味着接收方面的流量超出了硬件处理的能力。
	RxFifoErrors int64 `json:"rx_fifo_errors"` // 因FIFO缓冲区溢出而丢失的数据包数量
	// tx_fifo_errors: 发送FIFO缓冲区错误次数。同样，如果这个值很高，可能意味着发送操作正在受到限制。
	TxFifoErrors int64 `json:"tx_fifo_errors"` // 发送FIFO缓冲区错误次数
	// rx_missed_errors: 表示硬件遗漏的数据包数量。这个指标可能也反映了接收方面的瓶颈。
	RxMissedErrors int64 `json:"rx_missed_errors"` // 硬件遗漏的数据包数量
	// tx_carrier_errors: 网络载波相关的错误次数。这些错误通常跟物理层级上的问题有关，如断线、坏连接等。
	TxCarrierErrors int64 `json:"tx_carrier_errors"` // 网络载波相关的错误次数
}


type Stats struct {
	NetStats
	BufferSizes
	InterfaceName  string   `json:"interface_name"`
	InterfaceNames []string `json:"interface_names"`
}

func GetStats(interfaceName string) (stats Stats, err error) {
	var err1, err2 error
	interfaces, err := net.Interfaces()
	if err != nil {
		return stats, err
	}
	if len(interfaces) == 0 {
		return stats, errors.New("no interface")
	}
	stats.InterfaceNames = lo.Map(interfaces, func(inf net.Interface, i int) string {
		return inf.Name
	})
	if interfaceName == "" {
		interfaceName, _ = getMainInterfaceName()
	}
	stats.InterfaceName = interfaceName

	stats.NetStats, err1 = getNetstats(interfaceName)
	stats.BufferSizes, err2 = getBufferStats()
	return stats, errors.Join(err1, err2)
}

func getMainInterfaceName() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	if len(interfaces) == 0 {
		return "", errors.New("no interface")
	}
	type strint struct{
		name string
		bytes int64
	}
	data:= []strint{}
	for _, inf := range interfaces {
		if inf.Flags&net.FlagUp != 0 && inf.Flags&net.FlagLoopback == 0 {
			return inf.Name, nil
		}
		rx, _ := readNetStat(inf.Name, "rx_bytes")
		tx, _ := readNetStat(inf.Name, "tx_bytes")
		data = append(data, strint{inf.Name, tx+rx})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].bytes > data[j].bytes
	})
	return data[0].name, nil
}

// readNetStat 从对应的文件中读取网络统计值
func readNetStat(interfaceName string, stat string) (int64, error) {
	data, err := os.ReadFile(filepath.Join("/sys/class/net", interfaceName, "statistics", stat))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}


func getNetstats(interfaceName string) (NetStats, error) {
	// interfaceName := "eth0" // 替换为您实际的网卡接口名称
	stats := NetStats{}
	var err error

	// 逐一读取统计数据并存储到结构体中
	stats.RxBytes, err = readNetStat(interfaceName, "rx_bytes")

	stats.TxBytes, _ = readNetStat(interfaceName, "tx_bytes")

	stats.RxErrors, _ = readNetStat(interfaceName, "rx_errors")

	stats.TxErrors, _ = readNetStat(interfaceName, "tx_errors")

	stats.RxDropped, _ = readNetStat(interfaceName, "rx_dropped")

	stats.TxDropped, _ = readNetStat(interfaceName, "tx_dropped")

	stats.RxFifoErrors, _ = readNetStat(interfaceName, "rx_fifo_errors")

	stats.TxFifoErrors, _ = readNetStat(interfaceName, "tx_fifo_errors")

	stats.RxMissedErrors, _ = readNetStat(interfaceName, "rx_missed_errors")

	stats.TxCarrierErrors, _ = readNetStat(interfaceName, "tx_carrier_errors")
	return stats, err

}

// BufferSizes 存储了缓冲区大小信息
type BufferSizes struct {
	TCPReadMem     string `json:"tcp_rmem"`     // TCP接收缓冲区大小bytes (最小值、默认值、最大值)
	TCPWriteMem    string `json:"tcp_wmem"`     // TCP发送缓冲区大小 (最小值、默认值、最大值)
	TCPMem     string `json:"tcp_mem"`     // TCP缓冲区大小 (最小值、默认值、最大值)
	UDPMem     string `json:"udp_mem"`     // UDP缓冲区大小 (最小值、默认值、最大值)
	ReadMem     int64  `json:"rmem_min"` // 接收的最小缓冲区大小
	WriteMem    int64  `json:"wmem_min"` // 发送操作的最小缓冲区大小
	ReadMemMax  int64  `json:"rmem_max"` // 接收操作的最大缓冲区大小
	WriteMemMax int64  `json:"wmem_max"` // 发送操作的最大缓冲区大小
}

// readSysctl 从 /proc/sys 目录中读取指定的内核参数
func readSysctl(param string) (string, error) {
	data, err := os.ReadFile(filepath.Join("/proc/sys", strings.Replace(param, ".", "/", -1)))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func strconvToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func getBufferStats() (bufferSizes BufferSizes, err error) {
	// 获取TCP缓冲区大小(优化级最高): 最小值 默认值 最大值(bytes)
	bufferSizes.TCPReadMem, _ = readSysctl("net.ipv4.tcp_rmem")
	bufferSizes.TCPWriteMem, _ = readSysctl("net.ipv4.tcp_wmem")

	// 获取TCP缓冲区大小(tcp全局): 最小值 默认值 最大值(bytes)
	bufferSizes.TCPMem, _ = readSysctl("net.ipv4.tcp_mem")
	// 获取UDP缓冲区大小((tcp全局)): 最小值 默认值 最大值(bytes)
	bufferSizes.UDPMem, _ = readSysctl("net.ipv4.udp_mem")

	// 获取缓冲区大小(系统全局)
	rmemMin, _ := readSysctl("net.core.rmem_default")
	bufferSizes.ReadMem, err = strconvToInt64(rmemMin)

	wmemMin, _ := readSysctl("net.core.wmem_default")

	bufferSizes.WriteMem, _ = strconvToInt64(wmemMin)

	rmemMax, _ := readSysctl("net.core.rmem_max")

	bufferSizes.ReadMemMax, _ = strconvToInt64(rmemMax)

	wmemMax, _ := readSysctl("net.core.wmem_max")

	bufferSizes.WriteMemMax, _ = strconvToInt64(wmemMax)
	return bufferSizes, err
}
