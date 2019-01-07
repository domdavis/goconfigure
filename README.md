# Configuration Package for Go

[![Build Status](https://travis-ci.org/domdavis/goconfigure.svg?branch=master)](https://travis-ci.org/domdavis/goconfigure)
[![Coverage Status](https://coveralls.io/repos/github/domdavis/goconfigure/badge.svg?branch=master)](https://coveralls.io/github/domdavis/goconfigure?branch=master)
[![](https://godoc.org/github.com/domdavis/goconfigure?status.svg)](http://godoc.org/github.com/domdavis/goconfigure)

`goconfigure` allows configuration of an application via flags, environment
variables, or a configuration file. 

## Installation

```
go get github.com/domdavis/goconfigure
```


## Usage

The following comes from `example/cli.go` which can be used to see how the
various methods of setting options works:

```go
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
	opt.EnvVar("GOCONFIGURE_COUNT")
	opt.Default(1)
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.config,
		"The path to the configuration file")
	opt.ShortFlag('o')
	opts.Add(opt)

	if err := opts.ParseUsing(opt); err != nil {
		opts.Usage()
	}

	for i := 0; i < options.count; i++ {
		fmt.Println(options.message)
	}
}
```

By default this produces the output:

```
$ go run ./cli.go 
This space intentionally left blank
```

Using an options file gives:

```
$ go run ./cli.go -o ./options.json
Hello, world!
Hello, world!
Hello, world!
```

And this can be overridden, with command line arguments always taking 
precedence:

```
$ GOCONFIGURE_COUNT=2 go run cli.go -o ./options.json 
Hello, world!
Hello, world!
```

```
$ GOCONFIGURE_COUNT=2 go run cli.go -o ./options.json --count 3
Hello, world!
Hello, world!
Hello, world!
```
