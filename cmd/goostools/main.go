package main

import (
	"flag"

	"github.com/ahuigo/goos-tools/helper"
	"github.com/ahuigo/goos-tools/nets"
)

func main() {
	interfaceNameP := flag.String("i", "", "network interface name")

    // Parse the command-line from os.Args[1:]. Must be called after all flags are defined and before flags are accessed by the program.
    flag.Parse()
	interfaceName := ""
	if interfaceNameP != nil {
		interfaceName = *interfaceNameP
	}

	stats, err := nets.GetStats(interfaceName)
	if err!=nil{
		panic(err)
	}
	helper.PrintPretty(stats)
}
