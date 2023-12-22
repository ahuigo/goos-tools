package helper

import (
	"encoding/json"
	"fmt"
)

func PrintPretty(d any){
	b, _ := json.MarshalIndent(d, "", "  ")
	fmt.Println(string(b))
}