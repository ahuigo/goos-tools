package nets

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func GetSockStat() (*SockStat, error) {
	stat := &SockStat{}
	stat.PageSize = os.Getpagesize()
	file := "/proc/net/sockstat"
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch fields[0] {
		case "sockets:":
			stat.Sockets, _ = strconv.Atoi(fields[2])
		case "TCP:":
			stat.TCP.InUse, _ = strconv.Atoi(fields[2])
			stat.TCP.Orphan, _ = strconv.Atoi(fields[4])
			stat.TCP.Tw, _ = strconv.Atoi(fields[6])
			stat.TCP.Alloc, _ = strconv.Atoi(fields[8])
			stat.TCP.Mem, _ = strconv.Atoi(fields[10])
		case "UDP:":
			stat.UDP.InUse, _ = strconv.Atoi(fields[2])
			stat.UDP.Mem, _ = strconv.Atoi(fields[4])
		case "UDPLITE:":
			stat.UDPLITE.InUse, _ = strconv.Atoi(fields[2])
		case "RAW:":
			stat.RAW.InUse, _ = strconv.Atoi(fields[2])
		case "FRAG:":
			stat.FRAG.InUse, _ = strconv.Atoi(fields[2])
			stat.FRAG.Memory, _ = strconv.Atoi(fields[4])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if stat6, err := getSockStat6(); err != nil {
		return nil, err
	} else {
		stat.TCP.InUse += stat6.TCP.InUse
		stat.UDP.InUse += stat6.UDP.InUse
		stat.UDPLITE.InUse += stat6.UDPLITE.InUse
		stat.RAW.InUse += stat6.RAW.InUse
		stat.FRAG.InUse += stat6.FRAG.InUse
		stat.FRAG.Memory += stat6.FRAG.Memory
	}

	return stat, nil
}

/*
bash# cat /proc/net/sockstat6
TCP6: inuse 78
UDP6: inuse 0
UDPLITE6: inuse 0
RAW6: inuse 0
FRAG6: inuse 0 memory 0
*/
func getSockStat6() (*SockStat, error) {
	file := "/proc/net/sockstat6"
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat := &SockStat{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch fields[0] {
		case "TCP6:":
			stat.TCP.InUse, _ = strconv.Atoi(fields[2])
		case "UDP6:":
			stat.UDP.InUse, _ = strconv.Atoi(fields[2])
		case "UDPLITE6:":
			stat.UDPLITE.InUse, _ = strconv.Atoi(fields[2])
		case "RAW6:":
			stat.RAW.InUse, _ = strconv.Atoi(fields[2])
		case "FRAG6:":
			stat.FRAG.InUse, _ = strconv.Atoi(fields[2])
			stat.FRAG.Memory, _ = strconv.Atoi(fields[4])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stat, nil
}
