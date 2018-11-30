package lib

import (
	"fmt"
	"os"
	"time"
	"log"
	"bufio"
	"context"
	"sync"
	"net"
	"strings"
	"github.com/fengdingbo/sub-domain-scanner/lib/dns"
)

type Result struct {
	Host string
	Addr []net.IP
}

type Scanner struct {
	opts       *Options
	resultChan chan Result
	wordChan   chan string
	count      int
	issued     int
	context    context.Context
	log        *os.File
	mu         *sync.RWMutex
	timeStart  time.Time
	BlackIPs   []net.IP
}

func NewScanner(opts *Options) *Scanner {
	var this Scanner
	this.opts = opts

	this.wordChan = make(chan string, this.opts.Threads)
	this.resultChan = make(chan Result)
	this.mu = new(sync.RWMutex)
	this.BlackIPs = []net.IP{}

	f, err := os.Create(opts.Log)
	this.log = f

	if err != nil {
		log.Fatalln(err)
	}
	return &this
}
func (this *Scanner) WildcardsDomain() {
	log.Printf("[+] Validate wildcard domain *.%v exists", this.opts.Domain)
	if ip, ok := this.IsWildcardsDomain(); ok {
		log.Printf("Domain %v is wildcard,*.%v ip is %v", this.opts.Domain, this.opts.Domain, ip)
		if ! this.opts.WildcardDomain {
			os.Exit(0)
		}

		for _, v := range ip {
			this.BlackIPs = append(this.BlackIPs, v)
		}
	}
}

func (this *Scanner) Start() {
	this.WildcardsDomain()

	this.timeStart = time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go this.progressPrint(&wg)

	for i := 0; i < this.opts.Threads; i++ {
		go this.worker(&wg)
	}

	// 读取字典
	f, err := os.Open(this.opts.Dict)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	this.count = getCountLine(f)

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		wg.Add(1)
		word := strings.TrimSpace(scanner.Text())
		this.wordChan <- word
	}

	wg.Wait()

	format := "\r%d|%d|%.4f%%|%.4f/s|scanned in %.2f seconds\n"
	this.mu.RLock()
	this.progressClean()
	fmt.Fprintf(os.Stderr, format,
		this.issued,
		this.count,
		float64(this.issued)/float64(this.count)*100,
		float64(this.issued)/time.Since(this.timeStart).Seconds(),
		time.Since(this.timeStart).Seconds(),
	)
	this.mu.RUnlock()

}

func (this *Scanner) incr() {
	this.mu.Lock()
	this.issued++
	this.mu.Unlock()
}

func (this *Scanner) worker(wg *sync.WaitGroup) {
	for v := range this.wordChan {
		this.incr()
		host := fmt.Sprintf("%s.%s", v, this.opts.Domain)
		ip, err := this.LookupHost(host)
		if err == nil {
			//this.resultChan <- Result{host, ip}
			this.result(Result{host, ip})
		}

		wg.Done()
	}
}

func (this *Scanner) result(re Result) {
	// 如果没有一个可用ip存在,则不记录
	if this.IsBlackIPs(re.Addr) {
		return
	}

	this.progressClean()
	fmt.Printf("[+] %v\n", re)

	this.mu.Lock()
	this.log.WriteString(fmt.Sprintf("%v\t%v\n", re.Host, re.Addr))
	this.mu.Unlock()
}

func (this *Scanner) progressClean() {
	fmt.Fprint(os.Stderr, "\r\x1b[2K")
}

func (this *Scanner) progressPrint(wg *sync.WaitGroup) {
	//start := time.Now()
	tick := time.NewTicker(1 * time.Second)
	format := "\r%d|%d|%.4f%%|%.4f/s|scanned in %.2f seconds"
	//log.Println("Starting")

Loop:
	for {
		select {
		case <-tick.C:
			this.mu.RLock()
			fmt.Fprintf(os.Stderr, format,
				this.issued,
				this.count,
				float64(this.issued)/float64(this.count)*100,
				float64(this.issued)/time.Since(this.timeStart).Seconds(),
				time.Since(this.timeStart).Seconds(),
			)
			this.mu.RUnlock()
			// Force quit
			if this.issued == this.count {
				break Loop;
			}
		}
	}

	wg.Done()
}

// 获取泛域名ip地址
func (this *Scanner) IsWildcardsDomain() (ip []net.IP, ok bool) {
	// Go package net exists bug?
	// @link https://github.com/golang/go/issues/28947
	// Nonsupport RFC 4592
	// net.LookupHost("*.qzone.qq.com") //  --> lookup *.qzone.qq.com: no such host

	// md5(random string)
	// byte := md5.Sum([]byte(time.Now().String()))
	// randSub:=hex.EncodeToString(byte[:])
	// host := fmt.Sprintf("%s.%s", randSub, this.opts.Domain)
	// addrs, err := net.LookupHost(host)

	host := fmt.Sprintf("*.%s", this.opts.Domain)
	addrs, err := this.LookupHost(host)

	if err != nil {
		return addrs, false
	}

	return addrs, true
}

func (this *Scanner) TestAXFR() (results []string, err error) {
	if results, err = dns.Axrf(this.opts.Domain); err == nil {
		for _, v := range results {
			this.log.WriteString(fmt.Sprintf("%s\n", v))
		}
	}
	return
}

func (this *Scanner) TestDNSServer() bool {
	ipaddr, err := this.LookupHost("google-public-dns-a.google.com") // test lookup an existed domain

	if err != nil {
		//log.Println(err)
		return false
	}
	// Validate dns pollution
	if ipaddr[0].String() != "8.8.8.8" {
		// Non-existed domain test
		_, err := this.LookupHost("test.bad.dns.fengdingbo.com")
		// Bad DNS Server
		if err == nil {
			return false
		}
	}

	return true
}

func getCountLine(f *os.File) int {
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
	}

	return i
}
