package gonic

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func getOOMScore() (score int, score_adj int) {
	pid := os.Getpid()
	score, _ = getFileNum(fmt.Sprintf("/proc/%d/oom_score", pid))
	score_adj, _ = getFileNum(fmt.Sprintf("/proc/%d/oom_score_adj", pid))
	return
}

func getFileNum(path string) (int, error) {
	b, _ := os.ReadFile(path)
	b = bytes.TrimSpace(b)
	return strconv.Atoi(string(b))
}
