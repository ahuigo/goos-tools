package gonic

import (
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

func getFileNum (path string) (int, error) {
	s, _ := os.ReadFile(path)
	return strconv.Atoi(string(s))
}