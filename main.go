package main

import (
	"os"
	"github.com/fengdingbo/sub-domain-scanner/lib"
	"flag"
	"log"
)

func loadOptions() *lib.Options {
	o := lib.New()
	flag.IntVar(&o.Threads, "t", 200, "Num of scan threads")
	flag.StringVar(&o.Domain, "d", "", "The target Domain")
	flag.StringVar(&o.Dict, "f", "dict/subnames_full.txt", "File contains new line delimited subs")
	flag.BoolVar(&o.Help, "h", false, "Show this help message and exit")
	flag.StringVar(&o.Log, "o", "", "Output file to write results to (defaults to ./log/{target}).txt")
	flag.StringVar(&o.DNSServer, "dns", "", "DNS global server,eg:8.8.8.8")
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

	// TODO 泛域名处理逻辑
	log.Printf("[+] Validate wildcard domain *.%v exists", o.Domain)
	if ip, ok := this.IsWildcardsDomain(); ok {
		log.Printf("Domain %v is wildcard,*.%v ip is %s", o.Domain, o.Domain, ip)
		os.Exit(0)
	}
	this.Start()

}
