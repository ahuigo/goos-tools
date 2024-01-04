package nets

import "errors"

func GetSockStat() (*SockStat, error) {
    stat := &SockStat{}
    return stat, errors.New("Not implemented")
}