package netstat

type TcpConnection struct {
	Proto string
	// Recv-Q 表示接收队列的长度，即等待应用程序读取的数据量。如果 Recv-Q 的值很大，说明应用程序没有及时读取数据，或者网络拥塞导致数据传输缓慢。
	RecvQ int

	// Send-Q 表示发送队列的长度，即等待发送的数据量。如果 Send-Q 的值很大，说明应用程序发送的数据量很大，或者网络拥塞导致数据传输缓慢。
	SendQ       int
	LocalAddr   string
	ForeignAddr string
	State       string
	Program     string
	Pid         string
}
