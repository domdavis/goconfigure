package goconfigure_test

import (
	"flag"
	"fmt"
	"github.com/domdavis/goconfigure"
	"os"
	"strings"
	"testing"
)

func ExampleOption_Default() {
	var options struct {
		test string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.ShortFlag('f')
	opt.LongFlag("flag")
	opt.Default("default value")
	opts.Add(opt)
	err := opts.Parse(nil)

	if err == nil {
		fmt.Println(options)
	} else {
		fmt.Println(err)
	}

	// Output:
	// {default value}
}

func ExampleOption_ConfigKey() {
	var options struct {
		test string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.ConfigKey("value")
	opt.Default("default value")
	opts.Add(opt)
	err := opts.Parse(map[string]interface{}{"value": "config value"})

	if err == nil {
		fmt.Println(options)
	} else {
		fmt.Println(err)
	}

	// Output:
	// {config value}
}

func ExampleOption_EnvVar() {
	const name = "GOCONFIGURE_TEST_VALUE"
	var options struct {
		test string
	}

	_ = os.Setenv(name, "test")
	opts := goconfigure.NewOptionsWithArgs([]string{})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.EnvVar(name)
	opts.Add(opt)
	err := opts.Parse(nil)

	if err == nil {
		fmt.Println(options)
	} else {
		fmt.Println(err)
	}

	// Output:
	// {test}
}

func ExampleOption_ShortFlag() {
	var options struct {
		test string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{"-f", "test"})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.ShortFlag('f')
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

func ExampleOption_Flags() {
	var options struct {
		test string
	}

	opts := goconfigure.NewOptionsWithArgs([]string{"-f", "test"})
	opt := goconfigure.NewOption(&options.test, "test option")
	opt.Flags('f', "flag")
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

func TestNewOption(t *testing.T) {
	t.Run("option.New(nil, string) will not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("option.New(nil, string) caused panic")
			}
		}()
		if opt := goconfigure.NewOption(nil, ""); opt == nil {
			t.Errorf("Failed to generate options with nil value")
		}
	})
}

func TestRegisterFlags(t *testing.T) {
	t.Run("A nil FlagSet will not error", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "")
		if err := opt.RegisterFlags(nil); err != nil {
			t.Errorf("unexpected error registering flags: %s", err)
		}
	})

	t.Run("An invalid option type will error (long flag)", func(t *testing.T) {
		var invalid []string
		opt := goconfigure.NewOption(invalid, "")
		opt.LongFlag("long")
		err := opt.RegisterFlags(flag.CommandLine)

		expected := "failed to set long flag: invalid option type for flag " +
			"\"long\": []string"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error registering flags: %v", err)
		}
	})

	t.Run("An invalid option type will error (short flag)", func(t *testing.T) {
		var invalid []string
		opt := goconfigure.NewOption(invalid, "")
		opt.ShortFlag('s')
		err := opt.RegisterFlags(flag.CommandLine)

		expected := "failed to set short flag: invalid option type for flag " +
			"\"s\": []string"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error registering flags: %v", err)
		}
	})

	t.Run("An invalid default will error", func(t *testing.T) {
		var value string
		opt := goconfigure.NewOption(&value, "")
		opt.Default([]string{})
		opt.ShortFlag('s')
		err := opt.RegisterFlags(flag.CommandLine)

		expected := "failed to set short flag: cannot use default option [] " +
			"([]string) as *string for flag s"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error registering flags: %v", err)
		}
	})
}

func TestOption_Parse(t *testing.T) {
	t.Run("Parsing an empty option will error", func(t *testing.T) {
		err := goconfigure.NewOption(nil, "bad").Parse(nil)
		expected := "option with description 'bad' not registered with a value"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing option: %v", err)
		}
	})

	t.Run("Parsing a non pointer option will error", func(t *testing.T) {
		err := goconfigure.NewOption("non-pointer", "bad").Parse(nil)
		expected := "type Option requires pointer value(*string, not string)"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing option: %v", err)
		}
	})

	t.Run("Parsing an incompatible option type will error", func(t *testing.T) {
		var value []string
		opt := goconfigure.NewOption(&value, "incompatible")
		opt.Default("text")
		err := opt.Parse(nil)

		expected := "failed to set option: value.Data string 'text', " +
			"failed to assign to type []string"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing option: %v", err)
		}
	})

	t.Run("Parsing incompatible config will error", func(t *testing.T) {
		var value string
		opt := goconfigure.NewOption(&value, "incompatible")
		opt.ConfigKey("key")
		err := opt.Parse(map[string]interface{}{"key":[]string{}})

		expected := "failed to parse option config: cannot convert config " +
			"type []string to string for 'key'"
		if err == nil || err.Error() != expected {
			t.Errorf("unexpected error parsing option: %v", err)
		}
	})

	t.Run("Parsing incompatible environments will error", func(t *testing.T) {
		const name = "GOCONFIGURE_TEST_VALUE"
		_ = os.Setenv(name, "text")
		var value int
		opt := goconfigure.NewOption(&value, "incompatible")
		opt.EnvVar(name)
		err := opt.Parse(nil)

		prefix := "failed to parse environment option " +
			"'GOCONFIGURE_TEST_VALUE': "
		if err == nil || !strings.HasPrefix(err.Error(), prefix) {
			t.Errorf("unexpected error parsing option: %v", err)
		}
	})
}

func TestOption_String(t *testing.T) {
	t.Run("No flags will be handled correctly", func(t *testing.T) {
		s := goconfigure.NewOption(nil, "An Example").String()

		if !strings.Contains(s, "No CLI option") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("Short and long flags will be shown", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.Flags('f', "flag")
		s := opt.String()

		if !strings.Contains(s, "-f, --flag") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("Short flags will be shown", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.ShortFlag('f')
		s := opt.String()

		if !strings.Contains(s, "-f") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("Long flags will be shown", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.LongFlag("flag")
		s := opt.String()

		if !strings.Contains(s, "--flag") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("String defaults will be quoted", func(t *testing.T) {
		opt := goconfigure.NewOption("", "An Example")
		opt.Default("default")
		s := opt.String()

		if !strings.Contains(s, "(default \"default\")") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("Non-string defaults will not be quoted", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.Default(1)
		s := opt.String()

		if !strings.Contains(s, "(default 1)") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})

	t.Run("Environment variables will be displayed", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.EnvVar("TEST_ENV")
		s := opt.String()

		if !strings.Contains(s, "TEST_ENV") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})
	
	t.Run("Config options will be displayed", func(t *testing.T) {
		opt := goconfigure.NewOption(nil, "An Example")
		opt.ConfigKey("key")
		s := opt.String()

		if !strings.Contains(s, "key") {
			t.Errorf("unexpected usage string:\n%s", s)
		}
	})
}
