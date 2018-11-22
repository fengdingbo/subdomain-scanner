package main

import (
	"os"
	"github.com/fengdingbo/sub-domain-scanner/lib"
	"flag"
	"fmt"
)

func main() {
	o:=lib.NewOptions()
	flag.IntVar(&o.Threads, "t", 50, "Num of scan threads")
	flag.StringVar(&o.Domain, "d", "", "The target Domain")
	flag.StringVar(&o.Wordlist, "w", "dict/subnames_full.txt", "Path to the wordlist")
	flag.BoolVar(&o.Help, "h", false, "Show this help message and exit")
	flag.StringVar(&o.Log, "o", "", "Output file to write results to (defaults to ./log/{target}).txt")
	flag.Parse()

	if !o.Validate() {
		flag.Usage()
		os.Exit(0)
	}

	ip,_:=o.GetExtensiveDomainIp()
	if ip != "" {
		fmt.Println("泛域名")
		os.Exit(0)
	}

	lib.Run(o)
}
