package lib

import (
	"net"
)

func (this *Scanner) IsBlackIPs(ips []net.IP) bool {
	i := len(ips);
	for _, v := range ips {
		if (this.IsBlackList(v.String())) {
			i--
		}
	}

	if i == 0 {
		return true
	}

	return false
}

func (this *Scanner) IsBlackList(s string) bool {
	if !IsPublicIP(net.ParseIP(s)) {
		return true
	}

	//blackIps := []string{"1.1.1.1", "127.0.0.1", "0.0.0.0", "0.0.0.1"}

	if _, ok := this.BlackList[s]; ok {
		return true
	}
	return false
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
