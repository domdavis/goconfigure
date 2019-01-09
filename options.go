package goconfigure

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Options holds a set of configuration options which can be provided by the
// command line, environment variables, configuration files, or default values.
type Options interface {
	// Add an option to this set of Options.
	Add(option Option)

	// Parse the Options using the provided map for configuration options.
	Parse(config map[string]interface{}) error

	// ParseUsing uses the given Option to locate and load the configuration
	// from a file.
	ParseUsing(option Option) error

	// NArg is the number of arguments remaining after flags have been
	// processed. Calling NArg before Parse will simply return 0.
	NArg() int

	// Args returns the non-flag command-line arguments. Calling Args before
	// calling Parse will simply return an empty slice.
	Args() []string

	// Usage displays the usage information for this set of options to STDERR.
	Usage()

	// UsageString building of custom usage output by providing just the usage
	//
	UsageString() string
}

type options struct {
	data  []Option
	args  []string
	flags *flag.FlagSet
}

// NewOptions returns a new Options type that takes its flags from the arguments
// provided to the process.
func NewOptions() Options {
	return NewOptionsWithArgs(os.Args[1:])
}

// NewOptionsWithArgs returns a new Options type that uses the given slice of
// strings as its argument set.
func NewOptionsWithArgs(args []string) Options {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.Usage = func() {}
	return &options{
		flags: flags,
		args:  args,
	}
}

func (o *options) Add(option Option) {
	o.data = append(o.data, option)
}

func (o *options) Parse(config map[string]interface{}) error {
	if err := o.parseFlags(); err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	if err := o.parseConfig(config); err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	return nil
}

func (o *options) ParseUsing(option Option) error {
	var file string
	config := map[string]interface{}{}

	if err := o.parseFlags(); err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	if option != nil {
		path := option.Value()
		if err := path.AssignTo(&file); err != nil {
			return fmt.Errorf("failed to read path for config file: %s", err)
		}
	}

	if file != "" {
		if b, err := ioutil.ReadFile(file); err != nil {
			return fmt.Errorf("error reading config %s: %s", file, err)
		} else if err := json.Unmarshal(b, &config); err != nil {
			return fmt.Errorf("error parsing config %s: %s", file, err)
		}
	}

	if err := o.parseConfig(config); err != nil {
		return fmt.Errorf("config error: %s", err)
	}

	return nil
}

func (o *options) NArg() int {
	return o.flags.NArg()
}

func (o *options) Args() []string {
	return o.flags.Args()
}

func (o *options) Usage() {
	b := strings.Builder{}
	b.WriteString("Usage of ")
	b.WriteString(os.Args[0])
	b.WriteString(":\n")
	b.WriteString(o.UsageString())
	_, _ = fmt.Fprint(os.Stderr, b.String())
}

func (o *options) UsageString() string {
	b := strings.Builder{}

	for _, opt := range o.data {
		b.WriteString(opt.String())
	}

	if len(o.data) == 0 {
		b.WriteString("    \tNo configuration options set")
	}

	b.WriteString("\n")
	return b.String()
}

func (o *options) parseFlags() error {
	for _, opt := range o.data {
		if err := opt.RegisterFlags(o.flags); err != nil {
			return fmt.Errorf("failed to register flags: %s", err)
		}
	}

	if err := o.flags.Parse(o.args); err != nil {
		return fmt.Errorf("failed to parse flags: %s", err)
	}

	return nil
}

func (o *options) parseConfig(config map[string]interface{}) error {
	for _, opt := range o.data {
		if err := opt.Parse(config); err != nil {
			return fmt.Errorf("error parsing options: %s", err)
		}
	}

	return nil
}
