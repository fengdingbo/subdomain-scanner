package main

import (
	"os"
	"github.com/fengdingbo/sub-domain-scanner/lib"
	"flag"
	"log"
)

func main() {
	o:=lib.New()
	flag.IntVar(&o.Threads, "t", 200, "Num of scan threads")
	flag.StringVar(&o.Domain, "d", "", "The target Domain")
	flag.StringVar(&o.Wordlist, "w", "dict/subnames_full.txt", "Dict to the wordlist")
	flag.BoolVar(&o.Help, "h", false, "Show this help message and exit")
	flag.StringVar(&o.Log, "o", "", "Output file to write results to (defaults to ./log/{target}).txt")
	flag.StringVar(&o.DNSAddress, "dns", "", "DNS global server,eg:8.8.8.8")
	flag.Parse()

	if !o.Validate() {
		flag.Usage()
		os.Exit(0)
	}


	log.Printf("[+] Validate DNS servers...")
	if !o.TestDNSServer() {
		log.Println("[!] DNS servers unreliable")
		os.Exit(0)
	}
	log.Printf("[+] Found DNS Server %s", o.DNSAddress)

	// TODO 泛域名处理逻辑
	log.Printf("[+] Validate extensive domain *.%v exists",o.Domain)
	if ip,ok:=o.GetExtensiveDomainIp();ok {
		log.Printf("Domain %v is extensive,*.%v ip is %s", o.Domain, o.Domain, ip)
		os.Exit(0)
	}

	lib.Run(o)
}
