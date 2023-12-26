package netstats

/**
$ cat /proc/net/sockstat
sockets: used 238
TCP: inuse 12 orphan 0 tw 0 alloc 15 mem 6
UDP: inuse 5 mem 2
UDPLITE: inuse 0
RAW: inuse 0
FRAG: inuse 0 memory 0
分析一下：
	sockets: used 241
		系统中当前使用的套接字总数。
	TCP: inuse 15 orphan 0 tw 0 alloc 18 mem 6
		inuse：当前打开的 TCP 套接字数量。
		orphan：没有关联进程的 TCP 套接字数量，通常是因为原来的进程已经结束，但是 TCP 连接还没有完全关闭。
		tw：处于 TIME-WAIT 状态的 TCP 套接字数量。TIME-WAIT 是 TCP 连接关闭过程中的一个状态。
		alloc：已分配但还未使用的 TCP 套接字数量。
		mem：TCP 缓冲区使用的内存页数(已使用)。默认，每个页的大小为 os.Getpagesize()=4KB
	UDP: inuse 5 mem 2
		inuse：当前打开的 UDP 套接字数量。
		mem：UDP 缓冲区使用的内存页数。
	UDPLITE: inuse 0 // UDPLite 是一种类似于 UDP 的协议，但是提供了可选的校验和功能。
	RAW: inuse 0	//原始套接字可以用于直接发送或接收 IP 数据包
	FRAG: inuse 0 memory 0
		inuse：当前正在处理的 IP 分片数量。
		memory：用于处理 IP 分片的内存页数。
*/
type SockStat struct {
    Sockets int `json:"sockets"`
    TCP     struct {
        InUse  int `json:"inuse"`
        Orphan int `json:"orphan"`
        Tw     int `json:"tw"`
        Alloc  int `json:"alloc"`
		// BTW: /proc/sys/net/ipv4/tcp_rmem 和 /proc/sys/net/ipv4/tcp_wmem 这两个文件分别表示 TCP 接收缓冲区和发送缓冲区的大小
        Mem    int `json:"mem"`
    } `json:"TCP"`
    UDP struct {
        InUse int `json:"inuse"`
        Mem   int `json:"mem"`
    } `json:"UDP"`
    UDPLITE struct {
        InUse int `json:"inuse"`
    } `json:"UDPLITE"`
    RAW struct {
        InUse int `json:"inuse"`
    } `json:"RAW"`
    FRAG struct {
        InUse  int `json:"inuse"`
        Memory int `json:"memory"`
    } `json:"FRAG"`
	PageSize int `json:"pagesize"` //4KB
}