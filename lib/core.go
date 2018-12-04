package lib

import (
	"fmt"
	"os"
	"time"
	"log"
	"bufio"
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
	opts      *Options          // Options
	wordChan  chan string       // 字典队列
	found     int               // 发现域名数
	count     int               // 字典总数
	issued    int               // 当前处理数
	log       *os.File          // log文件
	mu        *sync.RWMutex     // 锁
	timeStart time.Time         // 执行开始时间
	BlackList map[string]string // 解析的黑名单字典，ip->泛域名
}

func NewScanner(opts *Options) *Scanner {
	var this Scanner
	this.opts = opts

	this.wordChan = make(chan string, this.opts.Threads)
	this.mu = new(sync.RWMutex)
	this.BlackList = make(map[string]string)

	f, err := os.Create(opts.Log)
	this.log = f

	if err != nil {
		log.Fatalln(err)
	}
	return &this
}

// 验证泛域名
// 解析泛域名后将解析的结果加入黑名单
func (this *Scanner) WildcardsDomain(s string) bool {
	//log.Printf("[+] Validate wildcard domain *.%v exists", this.opts.Domain)
	if ip, ok := this.IsWildcardsDomain(s); ok {
		//log.Printf("[+] Domain %v is wildcard,*.%v ip is %v", this.opts.Domain, this.opts.Domain, ip)
		//if ! this.opts.WildcardDomain {
		//	return true
		//}

		for _, v := range ip {
			this.BlackList[v.String()] = fmt.Sprintf("*.%s", this.opts.Domain)
		}

		return true
	}

	return false
}

func (this *Scanner) Start() {
	if ok := this.WildcardsDomain(this.opts.Domain); ok && !this.opts.WildcardDomain {
		return
	}

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

	format := "All Done. %d found, %.4f/s, %d scanned in %.2f seconds\n"
	this.mu.RLock()
	this.progressClean()
	log.Printf(format,
		this.found,
		float64(this.issued)/time.Since(this.timeStart).Seconds(),
		this.issued,
		time.Since(this.timeStart).Seconds(),
	)
	log.Printf("The output result file is %s\n", this.opts.Log)
	this.mu.RUnlock()

}

func (this *Scanner) incr() {
	this.mu.Lock()
	this.issued++
	this.mu.Unlock()
}

// goroutine >= 1
// 负责扫描字典队列
func (this *Scanner) worker(wg *sync.WaitGroup) {
	for v := range this.wordChan {
		this.incr()
		host := fmt.Sprintf("%s.%s", v, this.opts.Domain)
		ip, err := this.LookupHost(host)
		if err == nil {
			//this.resultChan <- Result{host, ip}
			this.result(Result{host, ip}, wg)
		}

		wg.Done()
	}
}
func (this *Scanner) addChan(re Result, wg *sync.WaitGroup) {
	if ok := this.WildcardsDomain(re.Host); ok && !this.opts.WildcardDomain {
		return
	}

	// TODO
	f, err := os.Open("dict/next_sub_full.txt")
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		this.mu.Lock()
		this.count++
		this.mu.Unlock()

		wg.Add(1)
		word := strings.TrimSpace(scanner.Text())
		word = fmt.Sprintf("%s.%s", word, re.Host[0:strings.LastIndex(re.Host, this.opts.Domain)-1])
		this.wordChan <- word
	}

	defer f.Close()
}

func (this *Scanner) result(re Result, wg *sync.WaitGroup) {
	// 如果没有一个可用ip存在,则不记录
	if this.IsBlackIPs(re.Addr) {
		return
	}

	this.progressClean()
	go this.addChan(re, wg)

	//php.Sleep(1)
	fmt.Printf("[+] %v\n", re)

	this.mu.Lock()
	this.found++
	this.log.WriteString(fmt.Sprintf("%v\t%v\n", re.Host, re.Addr))
	this.mu.Unlock()
}

// 清除光标所在行的所有字符
func (this *Scanner) progressClean() {
	fmt.Fprint(os.Stderr, "\r\x1b[2K")
}

// goroutine = 1
// 启动后该方法负责打印进度
// 直到进度到100%跳出死循环
func (this *Scanner) progressPrint(wg *sync.WaitGroup) {
	tick := time.NewTicker(1 * time.Second)
	format := "\r%d|%.4f%%|%.4f/s|%d scanned in %.2f seconds"
	log.Println("Starting")

Loop:
	for {
		select {
		case <-tick.C:
			this.mu.RLock()
			fmt.Fprintf(os.Stderr, format,
				this.count,
				float64(this.issued)/float64(this.count)*100,
				float64(this.issued)/time.Since(this.timeStart).Seconds(),
				this.issued,
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
func (this *Scanner) IsWildcardsDomain(s string) (ip []net.IP, ok bool) {
	// Go package net exists bug?
	// @link https://github.com/golang/go/issues/28947
	// Nonsupport RFC 4592
	// net.LookupHost("*.qzone.qq.com") //  --> lookup *.qzone.qq.com: no such host

	// md5(random string)
	// byte := md5.Sum([]byte(time.Now().String()))
	// randSub:=hex.EncodeToString(byte[:])
	// host := fmt.Sprintf("%s.%s", randSub, this.opts.Domain)
	// addrs, err := net.LookupHost(host)

	host := fmt.Sprintf("*.%s", s)
	addrs, err := this.LookupHost(host)

	if err != nil {
		return addrs, false
	}

	return addrs, true
}

// 验证DNS域传送
func (this *Scanner) TestAXFR(domain string) (results []string, err error) {
	server, err := this.LookupNS(domain)

	if results, err = dns.Axrf(domain, server); err == nil {
		for _, v := range results {
			this.mu.Lock()
			this.log.WriteString(fmt.Sprintf("%s\n", v))
			this.mu.Unlock()
		}
	}
	return
}

// 验证DNS服务器是否稳定
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
