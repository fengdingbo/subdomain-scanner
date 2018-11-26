package lib

import (
	"net"
	"context"
	"time"
)

func (opts *Options)  DNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", opts.DNSAddress)
}

func (opts *Options) Dns(subDomain string,ch chan<- Result) {
	if subDomain=="" {
		ch<- Result{}
		return
	}

	host:= subDomain+"."+opts.Domain
	addrs, err:=opts.LookupHost(host)
	if err != nil {
		//fmt.Println(err)
		ch<- Result{}
		return
	}
	ch<- Result{Host:host, Addr:addrs}
	return
}

func  (opts *Options) LookupHost(host string) (addrs []string, err error) {
	r := net.Resolver{
		PreferGo:true,
		Dial: opts.DNSDialer,
	}

	ctx := context.Background()

	ctx,_=context.WithTimeout(ctx,100*time.Millisecond)
	ipaddr, err := r.LookupHost(ctx, host)
	if err != nil {
		return
	}

	return ipaddr,nil
}