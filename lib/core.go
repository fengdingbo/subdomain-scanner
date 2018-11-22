package lib

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"time"
	"net"
	"log"
)

type Result struct {
	Host string
	Addr []string
}

func (opts *Options) Dns(subDomain string,ch chan<- Result) {
	host:= subDomain+"."+opts.Domain
	ips, err := net.LookupHost(host)
	if err != nil {

		ch<- struct {
			Host string
			Addr []string
		}{Host: host, Addr: []string{}}
		//ch<-fmt.Sprint(err)
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
	fmt.Println("read dict...")

	// 创建空线程
	ch := make(chan Result)
	for i := 0; i <= opts.Threads; i++ {
		go opts.Dns("", ch)
	}


	output, err := os.Create("log/"+opts.Domain+ ".txt")
	if err != nil {
		log.Fatalf("error on creating output file: %v", err)
	}

	g:=0
	for scanner.Scan() {
		g++
		select {
		case re := <-ch:
			// 处理完一个，马上再添加一个
			// 线程添加，直到某结果集处理完
			go opts.Dns(strings.TrimSpace(scanner.Text()), ch)
			if len(re.Addr) > 0 {
				opts.resultWorker(output, re)
			}

			fmt.Printf("%d\r",g)
		case <-time.After(3 * time.Second):
			fmt.Println("3秒超时")
			//	os.Exit(2)
		}
	}

	fmt.Println("")
}

func (opts *Options) resultWorker(f *os.File, re Result) {
	// 如果没有一个可用ip存在,则不记录
	i:=len(re.Addr);
	for _,v:= range re.Addr{
		if (IsBlackIP(v)) {
			i--
		}
	}
	if i==0{
		return
	}

	code,_:=Head("http://"+re.Host)
	fmt.Println(re,code)

	writeToFile(f, fmt.Sprintf("%v\t%v",re.Host,re.Addr))
}

// 获取泛域名ip地址
func (opts *Options) GetExtensiveDomainIp() (ip string,err error)  {
	host := "*." + opts.Domain
	ns, err := net.LookupHost(host)

	if err != nil {
		//fmt.Fprintf(os.Stderr, "Err: %s\n", err.Error())
		return
	}

	ip = ns[0]

	return ip,nil
}


func writeToFile(f *os.File, output string) error {
	_, err := f.WriteString(fmt.Sprintf("%s\n", output))
	if err != nil {
		return fmt.Errorf("[!] Unable to write to file %v", err)
	}
	return nil
}

func Run(opts *Options) {
	opts.Start()
}