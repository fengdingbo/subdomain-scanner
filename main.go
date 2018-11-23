package main

import (
	"os"
	"github.com/fengdingbo/sub-domain-scanner/lib"
	"flag"
	"log"
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

	// TODO 泛域名处理逻辑
	log.Printf("Check Domain *.%v exists",o.Domain)
	ip,_:=o.GetExtensiveDomainIp()
	if ip != "" {
		log.Printf("Domain %v is extensive,*.%v ip is %s", o.Domain,o.Domain, ip)
		os.Exit(0)
	}

	lib.Run(o)
}
