package lib

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"time"
	"log"
	"strconv"
	"crypto/md5"
	"encoding/hex"
)

type Result struct {
	Host string
	Addr []string
}


// 获取泛域名ip地址
func (opts *Options) GetExtensiveDomainIp() (ip string,ok bool)  {
	// Go package net exists bug?
	// Nonsupport RFC 4592
	// https://github.com/golang/go/issues/28947
	// opts.LookupHost("*.qzone.qq.com") //  --> lookup *.qzone.qq.com: no such host

	byte := md5.Sum([]byte(time.Now().String()))
	randSub:=hex.EncodeToString(byte[:])

	host := fmt.Sprintf("%s.%s", randSub, opts.Domain)
	addrs, err := opts.LookupHost(host)

	if err == nil {
		return addrs[0], true
	}

	return "", false
}

func  (opts *Options) TestDNSServer() bool {
	ipaddr, err := opts.LookupHost("google-public-dns-a.google.com") // test lookup an existed domain

	if err != nil {
		log.Println(err)
		return false
	}
	// Validate dns pollution
	if ipaddr[0] != "8.8.8.8" {
		// Non-existed domain test
		_, err := opts.LookupHost("test.bad.dns.fengdingbo.com")
		// Bad DNS Server
		if err == nil {
			return false
		}
	}

	return true
}

func (opts *Options) Start( ) {
	output, err := os.Create(opts.Log)
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

	// 读取字典
	f, err := os.Open(opts.Wordlist)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer output.Close()

	scanner := bufio.NewScanner(f)
	log.Printf("Read dict...")
	log.Printf("Found dict total %d.", count)
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
			fmt.Fprintf(os.Stderr, format, i,count,float64(i)/float64(count)*100)
		case <-time.After(60 * time.Second):
			log.Println("1秒超时")
			//	os.Exit(0)
		}
	}

	// bug 最后N个没有被接收
LOOP:
	for i := 0; i < opts.Threads; i++ {
		select {
		case re := <-ch:
			if len(re.Addr) > 0 {
				opts.resultWorker(output, re)
			}
		case <-time.After(60 * time.Second):
			log.Println("1秒超时...")
			break LOOP;
		}
	}


	log.Printf("Log file --> %s", opts.Log)
	log.Printf("Scan done")
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

	log.Printf("%v\t%v",re.Host,re.Addr)

	writeToFile(f, fmt.Sprintf("%v\t%v",re.Host,re.Addr))
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