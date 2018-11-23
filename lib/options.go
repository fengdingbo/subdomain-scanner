package lib

import (
	"os"
)

type Options struct {
	Threads		int
	Domain		string
	Wordlist	string
	Help 		bool
	Log			string
	DNSAddress			string
}


func NewOptions() *Options {
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
	_, err := os.Stat(opts.Wordlist)
	if err!= nil {
		return false
	}

	if opts.Log == "" {
		opts.Log = "log/"+opts.Domain+".txt"
	}

	if opts.DNSAddress == "" {
		opts.DNSAddress = "8.8.8.8:53"
	}

	return true
}

