package lib

import (
	"os"
	"fmt"
	"reflect"
	"bufio"
)

type Options struct {
	Threads        int
	Domain         string
	Dict           string
	Help           bool
	Log            string
	DNSServer      string
	WildcardDomain bool
	AXFC           bool
	ScanListFN     string
	ScanDomainList []string
}

func New() *Options {
	return &Options{
	}
}

func (opts *Options) existsDomain() bool {
	opts.ScanDomainList = []string{}
	for {
		if opts.ScanListFN != "" {
			f, err := os.Open(opts.ScanListFN)
			if err != nil {
				break;
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				opts.ScanDomainList = append(opts.ScanDomainList, scanner.Text())
			}
			f.Close()
		}

		break;
	}

	if (len(opts.ScanDomainList) > 0) {
		return true
	}
	if opts.Domain != "" {
		opts.ScanDomainList = append(opts.ScanDomainList, opts.Domain)
		return true
	}

	return false
}

func (opts *Options) Validate() bool {
	if opts.Help {
		return false
	}

	if (! opts.existsDomain()) {
		return false
	}

	if opts.Threads <= 0 {
		return false
	}

	_, err := os.Stat(opts.Dict)
	if err != nil {
		return false
	}

	if opts.Log == "" {
		logDir := "log"
		_, err := os.Stat(logDir)
		if err != nil {
			os.Mkdir(logDir, os.ModePerm)
		}
		opts.Log = fmt.Sprintf("%s/%s.txt", logDir, opts.Domain)
	}

	if opts.DNSServer == "" {
		//=============================================
		// 114 DNS		114.114.114.114	114.114.115.115
		// 阿里 AliDNS	223.5.5.5	223.6.6.6
		// 百度 BaiduDNS	180.76.76.76
		// DNSPod DNS+	119.29.29.29	182.254.116.116
		// CNNIC SDNS	1.2.4.8	210.2.4.8
		// oneDNS		117.50.11.11	117.50.22.22
		// DNS 派
		// 电信/移动/铁通	101.226.4.6	218.30.118.6
		// DNS 派 联通	123.125.81.6	140.207.198.6
		// Google DNS	8.8.8.8	8.8.4.4
		// IBM Quad9	9.9.9.9
		// OpenDNS		208.67.222.222	208.67.220.220
		// V2EX DNS		199.91.73.222	178.79.131.110
		//=============================================
		opts.DNSServer = "8.8.8.8/8.8.4.4"
	}

	opts.printOptions()
	return true
}

func (opts *Options) printOptions() {
	value := reflect.ValueOf(*opts)
	types := reflect.TypeOf(*opts)

	fmt.Println(`=============================================
sub-domain-scanner v0.3#dev
=============================================`)

	for i := 0; i < types.NumField(); i++ {
		if types.Field(i).Name[0] >= 65 && types.Field(i).Name[0] <= 90 {
			fmt.Printf("[+] %-15s: %v\n", types.Field(i).Name, value.Field(i).Interface())
		}
	}
	fmt.Println("=============================================")
}
