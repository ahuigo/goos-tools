package netstats

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func GetSockStat() (*SockStat, error) {
    stat := &SockStat{}
	stat.PageSize =os.Getpagesize()
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

    return stat, nil
}