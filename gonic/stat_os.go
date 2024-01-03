package gonic

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func getOsStat() (*gin.H, error) {
	res := &gin.H{}
	// memory
	memory, err := memory.Get()
	if err != nil {
		err = fmt.Errorf("mem:%s", err)
		return res, err
	}

	toG := func(u uint64) string {
		n := float64(u) / (1 << 30) // 1024^3
		return fmt.Sprintf("%.3fG", n)
	}

	//go memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// fmt.Printf("Alloc = %v MiB\n", m.Alloc / 1024 / 1024)

	// cpu
	cpuInfo, err := cpu.Get()
	if err != nil {
		err = fmt.Errorf("cpu:%s", err)
		return res, err
	}

	res = &gin.H{
		"memory(G bytes)": gin.H{
			"total":        toG(memory.Total),
			"used":         toG(memory.Used),
			"go:heapAlloc": toG(m.HeapAlloc),
			"go:heapInuse": toG(m.HeapInuse),
			"free":         toG(memory.Free),
			"cached":       toG(memory.Cached),
		},
		"cpu(G ticks)": gin.H{
			"total":  toG(cpuInfo.Total),
			"user":   toG(cpuInfo.User),
			"system": toG(cpuInfo.System),
			"idle":   toG(cpuInfo.Idle),
		},
	}
	return res, nil
}

func OsStat(c *gin.Context) {
	beforeStat, _ := getOsStat()
	res := gin.H{
		// "buildDate": conf.BuildDate,
		"cpu":    runtime.GOMAXPROCS(0),
		"before": beforeStat,
		"act":    c.Query("act"),
	}
	switch c.Query("act") {
	case "gc":
		runtime.GC()
		if c.Query("type") == "os" {
			debug.FreeOSMemory()
		}
	}
	if c.Query("act") != "" {
		afterStat, _ := getOsStat()
		res["after"] = afterStat
	}

	c.JSON(http.StatusOK, res)
}
