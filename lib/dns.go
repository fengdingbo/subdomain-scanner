package lib

import (
	"net"
	"context"
	"fmt"
	"github.com/fengdingbo/sub-domain-scanner/lib/dns"
)

func (this *Scanner)  DNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", this.opts.DNSAddress)
}


func (this *Scanner) Run(subDomain string,ch chan<- Result) {
	if subDomain=="" {
		ch<- Result{}
		return
	}
	host := fmt.Sprintf("%s.%s",subDomain,this.opts.Domain)
	addrs, err:=this.LookupHost(host)
	if err != nil {
		//fmt.Println(err)
		ch<- Result{}
		return
	}
	ch<- Result{Host:host, Addr:addrs}
	return
}

func  (this *Scanner) LookupHost(host string) (addrs []net.IP, err error) {
	ipaddr,err := dns.LookupHost(host)
	if err != nil {
		//log.Println(err)
		return
	}

	return ipaddr,nil
}