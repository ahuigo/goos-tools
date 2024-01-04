package nets

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readNetStat 从对应的文件中读取网络统计值
func readNetStat(interfaceName string, stat string) (int64, error) {
	data, err := os.ReadFile(filepath.Join("/sys/class/net", interfaceName, "statistics", stat))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

// readSysctl("net.ipv4.tcp_rmem")从 /proc/sys 目录中读取指定的内核参数
func readSysctl(param string) (string, error) {
	data, err := os.ReadFile(filepath.Join("/proc/sys", strings.Replace(param, ".", "/", -1)))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}