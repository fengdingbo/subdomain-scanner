package dns

import (
	"net"
	"strings"
	"errors"
	"github.com/miekg/dns"
)

func Axrf(hostname string) (results []string, err error) {
	results = []string{}
	domain := strings.ToLower(hostname)
	servers, err := net.LookupNS(domain)
	if err != nil {
		return
	}

	fqdn := dns.Fqdn(domain)
	for _, server := range servers {
		results = []string{}

		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server.Host, "53"))
		if err != nil {
			continue
		}

		for a := range answerChan {

			if a.Error != nil {
				continue
			}

			for _, rr := range a.RR {
				results = append(results, rr.String())
			}
		}

		return results, nil
	}
	return results, errors.New("Transfer failed")
}
