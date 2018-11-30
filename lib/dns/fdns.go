package dns

import (
	"net"
	"github.com/miekg/dns"
	"strings"
	"math/rand"
	"errors"
	"time"
)

func New(servers []string) *DnsResolver {
	for i := range servers {
		servers[i] = net.JoinHostPort(servers[i], "53")
	}

	return &DnsResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

type DnsResolver struct {
	Servers    []string
	RetryTimes int
	r          *rand.Rand
}

//
// func LookupHost(domain string) ([]net.IP, error) {}

func (r *DnsResolver) LookupHost(host string) ([]net.IP, error) {
	return r.lookupHost(host, r.RetryTimes)
}

func (r *DnsResolver) lookupHost(host string, triesLeft int) ([]net.IP, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{dns.Fqdn(host), dns.TypeA, dns.ClassINET}
	in, err := dns.Exchange(m1, r.Servers[r.r.Intn(len(r.Servers))])

	result := []net.IP{}

	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") && triesLeft > 0 {
			triesLeft--
			return r.lookupHost(host, triesLeft)
		}
		return result, err
	}

	if in != nil && in.Rcode != dns.RcodeSuccess {
		return result, errors.New(dns.RcodeToString[in.Rcode])

	}

	if len(in.Answer) == 0 {
		return result, errors.New("Unknown")
	}

	for _, record := range in.Answer {
		if t, ok := record.(*dns.A); ok {
			result = append(result, t.A)
		}
	}
	return result, err
}

func (r *DnsResolver) LookupNS(host string) ([]string, error) {
	return r.lookupNS(host, r.RetryTimes)
}

func (r *DnsResolver) lookupNS(host string, triesLeft int) ([]string, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{dns.Fqdn(host), dns.TypeNS, dns.ClassINET}
	in, err := dns.Exchange(m1, r.Servers[r.r.Intn(len(r.Servers))])

	result := []string{}
	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") && triesLeft > 0 {
			triesLeft--
			return r.lookupNS(host, triesLeft)
		}
		return result, err
	}

	if in != nil && in.Rcode != dns.RcodeSuccess {
		return result, errors.New(dns.RcodeToString[in.Rcode])
	}

	if len(in.Answer) == 0 {
		return result, errors.New("Unknown")
	}

	for _, record := range in.Answer {
		if t, ok := record.(*dns.NS); ok {
			result = append(result, t.Ns)
		}
	}
	return result, err
}
