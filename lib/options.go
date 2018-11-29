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
		logDir:="log"
		_, err := os.Stat(logDir)
		if err!= nil {
			os.Mkdir(logDir, os.ModePerm)
		}
		opts.Log = fmt.Sprintf("%s/%s.txt",logDir, opts.Domain)
	}

	if opts.DNSAddress == "" {
		opts.DNSAddress = "8.8.8.8:53"
	}

	return true
}

