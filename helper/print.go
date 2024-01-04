package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func PrintPretty(d any) {
	b, _ := json.MarshalIndent(d, "", "  ")
	fmt.Println(string(b))
}

func ParseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
