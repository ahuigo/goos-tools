package gonic

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type MemoryStat struct {
	Total         string
	Used          string
	Free          string
	Cached        string
	GoHeapAlloc   string //go heap使用到的内存
	GoHeapInuse   string //go heap向操作系统申请的内存(包括GoHeapAlloc, 已经被gc回收但未复用的内存)
	GoOOMScore    int    // 这个文件显示了一个进程被 OOM Killer 选中的得分。得分越高，进程被杀死的可能性越大。这个得分是根据进程的内存使用和运行时间等因素自动计算的，你不能直接修改这个文件
	GoOOMScoreAdj int    // 允许你调整一个进程的 OOM 得分。你可以写入一个从 -1000 到 1000 的值到这个文件，这个值将会被添加到进程的 OOM 得分上。例如，如果你写入 -500，那么这个进程的 OOM 得分将会减少 500，这将减少它被 OOM Killer 杀死的可能性。
}
type CpuStat struct {
	Total  string
	User   string
	System string
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

	oomScore, oomScoreAdj := getOOMScore()
	res = osStat{
		Memory: MemoryStat{
			Total:         toG(memory.Total),
			Used:          toG(memory.Used),
			Free:          toG(memory.Free),
			Cached:        toG(memory.Cached),
			GoHeapAlloc:   toG(m.HeapAlloc),
			GoHeapInuse:   toG(m.HeapInuse),
			GoOOMScore:    oomScore,
			GoOOMScoreAdj: oomScoreAdj,
		},
		Cpu: CpuStat{
			Total:  toG(cpuInfo.Total),
			User:   toG(cpuInfo.User),
			System: toG(cpuInfo.System),
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
