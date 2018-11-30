package main

import (
	"os"
	"github.com/fengdingbo/sub-domain-scanner/lib"
	"flag"
	"log"
	"fmt"
)

func loadOptions() *lib.Options {
	o := lib.New()
	flag.IntVar(&o.Threads, "t", 200, "Num of scan threads")
	flag.StringVar(&o.Domain, "d", "", "The target Domain")
	flag.StringVar(&o.Dict, "f", "dict/subnames_full.txt", "File contains new line delimited subs")
	flag.BoolVar(&o.Help, "h", false, "Show this help message and exit")
	flag.StringVar(&o.Log, "o", "", "Output file to write results to (defaults to ./log/{target}).txt")
	flag.StringVar(&o.DNSServer, "dns", "8.8.8.8/8.8.4.4", "DNS global server")
	flag.BoolVar(&o.WildcardDomain, "fw", false, "Force scan with wildcard domain")
	flag.BoolVar(&o.AXFC, "axfr", false, "DNS Zone Transfer Protocol (AXFR) of RFC 5936")
	flag.Parse()

	if !o.Validate() {
		flag.Usage()
		os.Exit(0)
	}
	return o
}

func main() {
	o := loadOptions()

	this := lib.NewScanner(o)

	log.Printf("[+] Validate DNS servers...")
	if !this.TestDNSServer() {
		log.Println("[!] DNS servers unreliable")
		os.Exit(0)
	}
	log.Printf("[+] Found DNS Server %s", o.DNSServer)

	// 检查是否存在DNS zone transfer
	log.Printf("[+] Validate AXFR of DNS zone transfer ")
	if axfr, err := this.TestAXFR(); err == nil {
		for _, v := range axfr {
			fmt.Println(v)
		}
		log.Printf("[+] Found DNS Server exists DNS zone transfer")
		os.Exit(0)
	}

	this.Start()

}
