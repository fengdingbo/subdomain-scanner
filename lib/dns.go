package lib

import (
	"net"
	"github.com/fengdingbo/subdomain-scanner/lib/dns"
	"strings"
)

func (this *Scanner) LookupHost(host string) (addrs []net.IP, err error) {
	DnsResolver := dns.New(strings.Split(this.opts.DNSServer, "/"))

	ipaddr, err := DnsResolver.LookupHost(host)
	if err != nil {
		//log.Println(err)
		return
	}

	return ipaddr, nil
}

func (this *Scanner) LookupNS(host string) ([]string, error) {
	DnsResolver := dns.New(strings.Split(this.opts.DNSServer, "/"))

	return DnsResolver.LookupNS(host)
}
