package gonic

import (
	"os"

	"github.com/ahuigo/goos-tools/netstat"
	"github.com/gin-gonic/gin"
)

// NetStat: show net stat information like shell `netstat`
func NetStat(c *gin.Context) {
	conns, err := netstat.GetAllTcpConnections()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	s := formatNetwork(conns)
	c.Data(200, "text/html", s)
	// c.String(200, s)
}

func formatNetwork(conns []netstat.TcpConnection) []byte {
	hostname, _:= os.Hostname()
	statisic := map[string]int{}
	for _, conn := range conns {
		statisic[conn.State]++
	}
	data := map[string]interface{}{
		"title":     "Netstat",
		"hostname":  hostname,
		"conns":    conns,
		"statisic": statisic,
	}

	s, _ := render("tpl/netstat.tmpl", data)
	return s
}
