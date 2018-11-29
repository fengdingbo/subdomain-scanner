package lib

import (
	"net"
	"github.com/fengdingbo/sub-domain-scanner/lib/dns"
	"strings"
)

func (this *Scanner) LookupHost(host string) (addrs []net.IP, err error) {
	DnsResolver:=dns.New(strings.Split(this.opts.DNSServer, "/"))

	ipaddr, err := DnsResolver.LookupHost(host)
	if err != nil {
		//log.Println(err)
		return
	}

	return ipaddr, nil
}
