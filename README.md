subdomain-scanner
======
使用Golang编写的子域名检测程序，特点就是快、快、快。

扫描速度依赖于网络环境。1Mb带宽，200个goroutine，稳定1700左右/s的扫描速度。

默认为谷歌的DNS服务器，可自行配置其它DNS。


## Dependencies ##
>go get github.com/miekg/dns

>go get github.com/hashicorp/go-multierror


## Building ##
	go get github.com/miekg/dns
	go get github.com/hashicorp/go-multierror
	go get github.com/fengdingbo/subdomain-scanner
	cd $GOPATH/src/github.com/fengdingbo/subdomain-scanner/
	make
	./subdomain-scanner -h


## Download from releases##
Download compiled binaries from [releases](https://github.com/fengdingbo/subdomain-scanner/releases)


## Usage ##
	Usage of ./subdomain-scanner -h
	  -axfr
			DNS Zone Transfer Protocol (AXFR) of RFC 5936 (default true)
	  -d string
			The target Domain
	  -depth int
			Scan sub domain depth. range[>=1] (default 1)
	  -dns string
			DNS global server (default "8.8.8.8/8.8.4.4")
	  -f string
			File contains new line delimited subs (default "dict/subnames_full.txt")
	  -fw
			Force scan with wildcard domain (default true)
	  -h	Show this help message and exit
	  -l string
			The target Domain in file
	  -o string
			Output file to write results to (defaults to ./log/{target}).txt
	  -t int
			Num of scan threads (default 200)


## Examples ##
        

## Change Log  
* [2018-12-03] 
	* 更好的参数调用提示
* [2018-12-01] 
	* 支持DNS域传送
	* 泛域名识别+扫描(泛域名得到的ip加入黑名单，继续爆破非黑名单ip)
* [2018-11-30] 
	* 重构并发逻辑
	* go官方的net包，不够完善，好多RFC都不支持，比如[RFC 4592](https://github.com/golang/go/issues/28947)，所以使用了一个第三方包来做dns解析，提升扫描效率。
* [2018-11-27] 
	* Demo雏形


## TODO ##
  - [x] 可选dns服务器
  - [x] 自定义字典
  - [x] 并发扫描
  - [x] 泛域名识别+扫描(泛域名得到的ip加入黑名单，继续爆破非黑名单ip)
  - [x] 支持DNS域传送
  - [x] 从文件中获取需要检测的域名
  - [ ] 支持DNS AAAA，ipv6检测
  - [x] 深度扫描(多级子域名检测)
  - [ ] 自定义导出格式、计划支持txt、json等
  - [x] 更友好的参数调用提示
  - [ ] 支持api接口调用


## Thanks ##
[https://github.com/miekg/dns](https://github.com/miekg/dns)

[https://github.com/OJ/gobuster](https://github.com/OJ/gobuster)

[https://github.com/binaryfigments/axfr](https://github.com/binaryfigments/axfr)

[https://github.com/lijiejie/subDomainsBrute](https://github.com/lijiejie/subDomainsBrute)
