package lib

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"time"
	"net"
	"log"
	"strconv"
)

type Result struct {
	Host string
	Addr []string
}

func (opts *Options) Dns(subDomain string,ch chan<- Result) {
	if subDomain=="" {
		ch<- Result{}
		return
	}

	host:= subDomain+"."+opts.Domain
	ips, err := net.LookupHost(host)
	if err != nil {
		ch<- Result{}
		return
	}

	ch<- Result{Host:host, Addr:ips}
}

func (opts *Options) Start( ) {
	// 读取字典
	f, err := os.Open(opts.Wordlist)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	log.Println("read dict...")

	output, err := os.Create("log/"+opts.Domain+ ".txt")
	if err != nil {
		log.Fatalf("error on creating output file: %v", err)
	}

	i:=0
	count:=opts.GetFileCountLine()
	width:=len(strconv.Itoa(count))
	format:=fmt.Sprintf("%%%dd|%%%dd|%%.4f%%%%\r",width,width)


	// 创建空线程
	if count < opts.Threads {
		opts.Threads=count
	}
	ch := make(chan Result)
	for i := 0; i < opts.Threads; i++ {
		go opts.Dns("", ch)
	}

	for scanner.Scan() {
		i++
		select {
		case re := <-ch:
			// 处理完一个，马上再添加一个
			// 线程添加，直到某结果集处理完
			go opts.Dns(strings.TrimSpace(scanner.Text()), ch)
			if len(re.Addr) > 0 {
				opts.resultWorker(output, re)
			}
			fmt.Printf(format,i,count,float64(i)/float64(count)*100)
		case <-time.After(3 * time.Second):
			log.Println("3秒超时")
			//	os.Exit(2)
		}
	}

	// bug 最后N个没有被接收
LOOP:
	for i := 0; i < opts.Threads; i++ {
		select {
			case re := <-ch:
				fmt.Println(re)
			case <-time.After(3 * time.Second):
				log.Println("3秒超时...")
				break LOOP;
		}
	}


	log.Println("结束")
}

func (opts *Options) resultWorker(f *os.File, re Result) {
	// 如果没有一个可用ip存在,则不记录
	i:=len(re.Addr);
	for _,v:= range re.Addr{
		if (IsBlackIP(v)) {
			i--
		}
	}
	if i==0 {
		return
	}

	log.Println(re)

	writeToFile(f, fmt.Sprintf("%v\t%v",re.Host,re.Addr))
}

// 获取泛域名ip地址
func (opts *Options) GetExtensiveDomainIp() (ip string,err error)  {
	host := "*." + opts.Domain
	ns, err := net.LookupHost(host)

	if err != nil {
		return
	}

	ip = ns[0]

	return ip,nil
}


func writeToFile(f *os.File, output string) (err error) {
	_, err = f.WriteString(fmt.Sprintf("%s\n", output))
	if err != nil {
		return
	}
	return nil
}


func (opts *Options) GetFileCountLine() int {
	// 读取字典
	f, err := os.Open(opts.Wordlist)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	g:=0
	for scanner.Scan() {
		g++
	}

	return g
}

func Run(opts *Options) {
	opts.Start()
}