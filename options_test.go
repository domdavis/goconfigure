package goconfigure_test

import (
	"fmt"
	"github.com/domdavis/goconfigure"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func ExampleNewOptions() {
	var r string
	opts := goconfigure.NewOptions()

	// We need to replace opts so arguments sent be testing harnesses are not
	// included. Ordinarily this line can be omitted.
	opts = goconfigure.NewOptionsWithArgs(nil)

	opt := goconfigure.NewOption(&r, "option")
	opt.LongFlag("config")
	opts.Add(opt)
	err := opts.ParseUsing(opt)

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleNewOptionsWithArgs() {
	var options struct {
		test string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{"-flag", "test"})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.ShortFlag('f')
	opt.LongFlag("flag")
	opts.Add(opt)
	err := opts.ParseUsing(nil)

	if err == nil {
		fmt.Println(options)
	} else {
		fmt.Println(err)
	}

	// Output:
	// {test}
}

func ExampleOptions_ParseUsing() {
	var options struct {
		config     string
		numeric    int
		text       string
		overridden string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{
		"-f", "testdata/config.json",
		"-o", "overridden",
	})

	opt := goconfigure.NewOption(&options.text, "from environment")
	opt.EnvVar("UNSET_ENV_VAR")
	opt.Default("unset")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.numeric, "from config")
	opt.ConfigKey("numeric")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.overridden, "overridden value")
	opt.ConfigKey("overridden")
	opt.ShortFlag('o')
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.config, "config file")
	opt.ShortFlag('f')
	opts.Add(opt)

	if err := opts.ParseUsing(opt); err != nil {
		fmt.Println(err)
	}

	fmt.Println(options)

	// Output:
	// {testdata/config.json 4 unset overridden}
}

func ExampleOptions_Add() {
	var opt goconfigure.Option
	var options struct {
		boolean      bool
		integer      int
		long         int64
		unsigned     uint
		unsignedLong uint64
		float        float64
		text         string
		duration     time.Duration
	}

	opts := goconfigure.NewOptionsWithArgs([]string{
		"-b", "-i", "1", "-l", "2", "-u", "3", "-unsignedLong", "4",
		"-f", "5.6", "-t", "words", "-d", "60s",
	})

	opt = goconfigure.NewOption(&options.boolean, "boolean")
	opt.Flags('b', "boolean")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.integer, "integer")
	opt.Flags('i', "integer")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.long, "long")
	opt.Flags('l', "long")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.unsigned, "unsigned")
	opt.Flags('u', "unsigned")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.unsignedLong, "unsignedLong")
	opt.Flags('z', "unsignedLong")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.float, "float")
	opt.Flags('f', "float")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.text, "text")
	opt.Flags('t', "text")
	opts.Add(opt)

	opt = goconfigure.NewOption(&options.duration, "duration")
	opt.Flags('d', "duration")
	opts.Add(opt)

	err := opts.Parse(nil)

	if err == nil {
		fmt.Println(options)
	} else {
		fmt.Println(err)
	}

	// Output:
	// {true 1 2 3 4 5.6 words 60000000000}
}

func TestOptions_Parse(t *testing.T) {
	t.Run("Parsing with invalid flags will error", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs([]string{"--undefined"})
		err := opts.Parse(nil)

		expected := "config error: failed to parse flags: " +
			"flag provided but not defined: -undefined"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing flags with an invalid type will error", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs(nil)
		opt := goconfigure.NewOption([]string{}, "invalid")
		opt.ShortFlag('s')
		opts.Add(opt)
		err := opts.Parse(nil)

		expected := "config error: failed to register flags: failed to set " +
			"short flag: invalid option type for flag \"s\": []string"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing config with an invalid type will error", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs(nil)
		opt := goconfigure.NewOption([]string{}, "invalid")
		opt.ConfigKey("key")
		opts.Add(opt)
		err := opts.Parse(nil)

		expected := "config error: error parsing options: type Option " +
			"requires pointer value(*[]string, not []string)"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})
}

func TestOptions_ParseUsing(t *testing.T) {
	t.Run("Parsing with invalid flags will error", func(t *testing.T) {
		var config string

		opts := goconfigure.NewOptionsWithArgs([]string{"--undefined"})
		opt := goconfigure.NewOption(&config, "config")
		err := opts.ParseUsing(opt)

		expected := "config error: failed to parse flags: " +
			"flag provided but not defined: -undefined"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing an invalid config file will error", func(t *testing.T) {
		var config string

		opts := goconfigure.NewOptionsWithArgs([]string{"--config", "invalid"})
		opt := goconfigure.NewOption(&config, "config")
		opt.LongFlag("config")
		opts.Add(opt)
		err := opts.ParseUsing(opt)

		expected := "error reading config invalid: open invalid: no such " +
			"file or directory"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing with an invalid config will error", func(t *testing.T) {
		var config string

		path := filepath.Join("testdata", "invalid.json")
		opts := goconfigure.NewOptionsWithArgs([]string{
			"--config", path})
		opt := goconfigure.NewOption(&config, "config")
		opt.LongFlag("config")
		opts.Add(opt)
		err := opts.ParseUsing(opt)

		expected := fmt.Sprintf("error parsing config %s: json: cannot "+
			"unmarshal array into Go value of type map[string]interface {}",
			path)

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing with an incorrect config will error", func(t *testing.T) {
		var config struct {
			file  string
			value int
		}

		opts := goconfigure.NewOptionsWithArgs([]string{
			"--config", filepath.Join("testdata", "incorrect.json")})
		c := goconfigure.NewOption(&config.file, "config")
		c.LongFlag("config")
		opts.Add(c)

		v := goconfigure.NewOption(&config.value, "value")
		v.ConfigKey("key")
		opts.Add(v)
		err := opts.ParseUsing(c)

		expected := "config error: error parsing options: failed to parse " +
			"option config: cannot convert config type string to int for 'key'"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})

	t.Run("Parsing with an invalid option will error", func(t *testing.T) {
		var config int

		opts := goconfigure.NewOptionsWithArgs(nil)
		opt := goconfigure.NewOption(&config, "config")
		opt.Default(1)
		opts.Add(opt)
		err := opts.ParseUsing(opt)

		expected := "failed to read path for config file: value.Data: " +
			"invalid cast of int to string"

		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing options: %v", err)
		}
	})
}

func TestOptions_NArg(t *testing.T) {
	t.Run("NArgs doesn't fail if Parse hasn't been called", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs([]string{"command"})

		n := opts.NArg()

		if n != 0 {
			t.Errorf("Unexpected number of NArgs: %d", n)
		}
	})

	t.Run("Args doesn't fail if Parse hasn't been called", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs([]string{"command"})

		a := opts.Args()

		if len(a) != 0 {
			t.Errorf("Unexpected Args: %d", a)
		}
	})

	t.Run("NArgs returns the number of extra arguments", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs([]string{"command"})
		if err := opts.Parse(nil); err != nil {
			t.Errorf("Unxpted error parsing arguments: %s", err)
		}

		n := opts.NArg()

		if n != 1 {
			t.Errorf("Unexpected number of NArgs: %d", n)
		}
	})

	t.Run("Args returns the extra arguments", func(t *testing.T) {
		opts := goconfigure.NewOptionsWithArgs([]string{"command"})
		if err := opts.Parse(nil); err != nil {
			t.Errorf("Unxpted error parsing arguments: %s", err)
		}

		a := opts.Args()

		if len(a) != 1 || a[0] != "command" {
			t.Errorf("Unexpected Args: %s", a)
		}
	})
}

func TestOptions_Usage(t *testing.T) {
	t.Run("An empty options wont panic when running usage", func(t *testing.T) {
		opts := goconfigure.NewOptions()
		opts.Usage()
	})
}

func TestOptions_UsageString(t *testing.T) {
	t.Run("Usage handles no options", func(t *testing.T) {
		opts := goconfigure.NewOptions()
		s := opts.UsageString()

		if !strings.Contains(s, "No configuration options set") {
			t.Errorf("Unexpted usage output: %s", s)
		}
	})

	t.Run("Usage handles options", func(t *testing.T) {
		opts := goconfigure.NewOptions()
		opts.Add(goconfigure.NewOption(nil, "Test option"))
		s := opts.UsageString()

		if !strings.Contains(s, "Test option") {
			t.Errorf("Unexpted usage output: %s", s)
		}
	})
}
