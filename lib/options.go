package lib

import (
	"os"
	"fmt"
)

type Options struct {
	wordMap		[]string
	Threads		int
	Domain		string
	Dict		string
	Help 		bool
	Log			string
	DNSAddress	string
}


func New() *Options {
	return &Options{
	}
}

func (opts *Options) Validate() bool{
	if opts.Help {
		return false
	}

	if opts.Threads<=0 {
		return false
	}

	if opts.Domain=="" {
		return false
	}
	_, err := os.Stat(opts.Dict)
	if err!= nil {
		return false
	}

	if opts.Log == "" {
		opts.Log = fmt.Sprintf("%s.txt",opts.Domain)
	}

	if opts.DNSAddress == "" {
		opts.DNSAddress = "223.5.5.5:53"
	}

	return true
}

