sub-domain-scanner
======
A fast sub domain brute tool for pentesters.


## Dependencies ##
>go get github.com/miekg/dns


## Building ##
    go get github.com/miekg/dns
    go get github.com/fengdingbo/sub-domain-scanner
    cd $GOPATH/src/github.com/fengdingbo/sub-domain-scanner/
    make
    ./sub-domain-scanner -h

## Usage ##
        Usage of ./sub-domain-scanner:
          -axfr
                DNS Zone Transfer Protocol (AXFR) of RFC 5936 (default true)
          -d string
                The target Domain
          -dns string
                DNS global server (default "8.8.8.8/8.8.4.4")
          -f string
                File contains new line delimited subs (default "dict/subnames_full.txt")
          -fw
                Force scan with wildcard domain
          -h	Show this help message and exit
          -l string
                The target Domain in file
          -o string
                Output file to write results to (defaults to ./log/{target}).txt
          -t int
                Num of scan threads (default 200)

## Examples ##
        

## Change Log ##
* [2018-12-01] 
        * 泛域名识别+扫描(泛域名得到的ip加入黑名单，继续爆破非黑名单ip)
        * 支持DNS域传送
* [2018-11-30]
        * 重构并发逻辑
        * go官方的net包，不够完善，好多RFC都不支持，比如RFC 4592，所以使用了一个第三方包来做dns解析，提升扫描效率。
* [2018-11-27]
        * Demo雏形

## 子域名扫描功能描述 ##
  - [x] 可选dns服务器
  - [x] 自定义字典
  - [x] 并发扫描
  - [x] 泛域名识别+扫描(泛域名得到的ip加入黑名单，继续爆破非黑名单ip)
  - [x] 支持DNS域传送
  - [x] 从文件中获取需要检测的域名
  - [ ] 支持DNS AAAA，ipv6解析
  - [ ] 深度扫描(多级子域名)
  - [ ] 自定义导出格式、计划支持txt、json等
  - [ ] 更好的参数调用提示
  