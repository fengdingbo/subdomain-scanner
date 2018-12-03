sub-domain-scanner v0.3
======
A fast sub domain brute tool for pentesters.

官方的net包，不够完善，好多RFC都不支持，所以使用了一个第三方包。


## Dependencies ##
>go get github.com/miekg/dns


## 编译运行 ##
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


## 子域名扫描功能描述 ##
  - [x] 可选dns服务器
  - [x] 自定义字典
  - [x] 并发扫描
  - [x] 泛域名识别+扫描(泛域名得到的ip加入黑名单，继续爆破非黑名单ip)
  - [x] 支持DNS域传送
  - [x] 从文件中获取需要检测的域名
  - [ ] 深度扫描(多级子域名)
  - [ ] 自定义导出格式、计划支持txt、json等
  - [ ] 更好的参数调用提示
  