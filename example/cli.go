package main

import (
	"fmt"
	"github.com/domdavis/goconfigure"
)

func main() {
	var opt goconfigure.Option
	var options struct {
		config string
		message string
		count int
	}

	opts := goconfigure.NewOptions()

	opt = goconfigure.NewOption(&options.message, "The message to display")
	opt.ShortFlag('m')
	opt.ConfigKey("message")
	opt.Default("This space intentionally left blank")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.count,
		"The number of times to display the message")
	opt.LongFlag("count")
	opt.ConfigKey("count")
	opt.Default(1)
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.config,
		"The path to the configuration file")
	opt.ShortFlag('o')
	opt.Default("./options.json")
	opts.Add(opt)

	if err := opts.ParseUsing(opt); err != nil {
		opts.Usage()
	}

	for i := 0; i < options.count; i++ {
		fmt.Println(options.message)
	}
}
