package lib

import (
	"net"
	"github.com/fengdingbo/sub-domain-scanner/lib/dns"
)

func (this *Scanner) LookupHost(host string) (addrs []net.IP, err error) {
	ipaddr, err := dns.LookupHost(host)
	if err != nil {
		//log.Println(err)
		return
	}

	return ipaddr, nil
}
