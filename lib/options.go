package lib

type Options struct {
	Threads		int
	Domain		string
	Wordlist	string
	Help 		bool
	Log			string
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

	return true
}

