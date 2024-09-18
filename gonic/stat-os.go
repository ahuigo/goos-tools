package gonic

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type MemoryStat struct {
	Total       string
	Used        string
	Free        string
	Cached      string
	GoHeapAlloc string //go heap使用到的内存
	GoHeapInuse string //go heap向操作系统申请的内存(包括GoHeapAlloc, 已经被gc回收但未复用的内存)
	GoMallocs   string // 表示从程序启动到现在，已经分配的堆对象的数量
	GoFrees     string // 表示从程序启动到现在，已经释放的堆对象的数量
}
type CpuStat struct {
	Total  string
	User   string
	System string
	Goroutines int
	Idle   string
}

type osStat struct {
	Memory   MemoryStat
	Cpu      CpuStat
	Hostname string
}

func getOsStat() (res osStat, err error) {
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

	res = osStat{
		Memory: MemoryStat{
			Total:       toG(memory.Total),
			Used:        toG(memory.Used),
			Free:        toG(memory.Free),
			Cached:      toG(memory.Cached),
			GoHeapAlloc: toG(m.HeapAlloc),
			GoHeapInuse: toG(m.HeapInuse),
			GoMallocs:   strconv.FormatUint(m.Mallocs, 10),
			GoFrees:     strconv.FormatUint(m.Frees, 10),
		},
		Cpu: CpuStat{
			Total:  toG(cpuInfo.Total),
			User:   toG(cpuInfo.User),
			System: toG(cpuInfo.System),
			Goroutines:  runtime.NumGoroutine(), // Current number of goroutines
			Idle:   toG(cpuInfo.Idle),
		},
	}
	return res, nil
}

func OsStat(c *gin.Context) {
	osStat, _ := getOsStat()
	// _, isJson := c.Params.Get("json")
	_, isJson := c.Request.URL.Query()["json"]
	if !isJson {
		s, err := formatOs(osStat)
		if err != nil {
			c.String(400, err.Error())
			c.AbortWithError(400, err)
			return
		}
		c.Data(200, "text/html", s)
		return
	}
	res := gin.H{
		// "buildDate": conf.BuildDate,
		"cpu":    runtime.GOMAXPROCS(0),
		"before": osStat,
		"isJson": isJson,
		"act":    c.Query("act"),
	}

	// 2. execute gc action
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

func formatOs(osStat osStat) ([]byte, error) {
	hostname, _ := os.Hostname()
	data := gin.H{
		"hostname": hostname,
		"memory":   osStat.Memory,
		"cpu":      osStat.Cpu,
	}

	s, err := render("tpl/os-stat.tmpl", data)
	return s, err
}
