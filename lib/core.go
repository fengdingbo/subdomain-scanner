package lib

import (
	"fmt"
	"os"
	"time"
	"log"
	"bufio"
	"context"
	"strings"
	"sync"
	"net"
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
	mu         *sync.RWMutex
}

func NewScanner(opts *Options) *Scanner {
	var this Scanner
	this.opts = opts

	this.wordChan = make(chan string, this.opts.Threads)
	this.resultChan = make(chan Result)
	this.mu = new(sync.RWMutex)
	return &this
}

func (this *Scanner) Start() {
	for i := 0; i < this.opts.Threads; i++ {
		go this.worker()
	}

	go this.result()
	go this.progressPrint()

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
		word := strings.TrimSpace(scanner.Text())
		this.wordChan <- word
	}
	return
}

func (this *Scanner) incr() {
	this.mu.Lock()
	this.issued++
	this.mu.Unlock()
}

func (this *Scanner) worker() {
	for v := range this.wordChan {
		this.incr()
		host := fmt.Sprintf("%s.%s", v, this.opts.Domain)
		ip, err := this.LookupHost(host)
		if err == nil {
			this.resultChan <- Result{host, ip}
		}
	}

	fmt.Println("worker")
}

func (this *Scanner) result() {
	f, err := os.Create(this.opts.Log)
	for re := range this.resultChan {
		// 如果没有一个可用ip存在,则不记录
		if IsBlackIPs(re.Addr) {
			continue
		}

		this.progressClean()
		fmt.Printf("[+] %v\n", re)
		if err == nil {
			f.WriteString(fmt.Sprintf("%v\t%v\n", re.Host, re.Addr))
		}
	}
}

func (this *Scanner) progressClean() {
	fmt.Fprint(os.Stderr, "\r\x1b[2K")
}

func (this *Scanner) progressPrint() {
	start := time.Now()
	tick := time.NewTicker(1 * time.Second)
	format := "\r%d|%d|%.4f%%|scanned in %.2f seconds"

	log.Println("Starting")
	for {
		select {
		case <-tick.C:
			this.mu.RLock()
			fmt.Fprintf(os.Stderr, format,
				this.issued,
				this.count,
				float64(this.issued)/float64(this.count)*100,
				time.Since(start).Seconds(),
			)
			this.mu.RUnlock()

			// Force quit
			if this.issued == this.count {
				//this.progressClean()
				fmt.Println("")
				log.Println("Finished")
				os.Exit(0)
			}
		}
	}
}

// 获取泛域名ip地址
func (this *Scanner) IsWildcardsDomain() (ip string, ok bool) {
	// Go package net exists bug?
	// Nonsupport RFC 4592
	// https://github.com/golang/go/issues/28947
	// net.LookupHost("*.qzone.qq.com") //  --> lookup *.qzone.qq.com: no such host

	// md5(random string)
	// byte := md5.Sum([]byte(time.Now().String()))
	// randSub:=hex.EncodeToString(byte[:])
	// host := fmt.Sprintf("%s.%s", randSub, this.opts.Domain)
	// addrs, err := net.LookupHost(host)

	host := fmt.Sprintf("*.%s", this.opts.Domain)
	addrs, err := this.LookupHost(host)

	if err == nil {
		return addrs[0].String(), true
	}

	return "", false
}

func (this *Scanner) TestDNSServer() bool {
	ipaddr, err := this.LookupHost("google-public-dns-a.google.com") // test lookup an existed domain

	if err != nil {
		log.Println(err)
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
