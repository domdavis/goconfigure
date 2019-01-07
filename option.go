package goconfigure

import (
	"flag"
	"fmt"
	"github.com/domdavis/goconfigure/value"
	"os"
	"reflect"
	"time"
)

// Option represents a configuration option that can be set either by flag,
// configuration file, environment variable, or a default value with the value
// to use being chosen in that order. Options must be one of bool, int, int64,
// uint, unit64, float64, string, or time.Duration.
type Option interface {

	// Flags defines both a short and long flag for setting the option from the
	// command line. For example:
	//
	//     option.Flags('f', "flag")
	//
	// allows the flags:
	//
	//     myApp -f value
	//     myApp --flag value
	//
	// An Option can only have one short and long flag definition so any
	// previously defined values will be overwritten. This function does nothing
	// if Parse has already been called.
	Flags(short rune, longFlag string)

	// ShortFlag defines a short flag for setting the option from the command
	// line. For example:
	//
	//     option.ShortFlag('f')
	//
	// allows the flag:
	//
	//     myApp -f value
	//
	// An Option can only have one short flag defined so any previously defined
	// value will be overwritten. This function does nothing if Parse has
	// already been called.
	ShortFlag(name rune)

	// LongFlag defines a long flag for setting the option from the command
	// line. For example:
	//
	//     option.LongFlag("flag")
	//
	// allows the flag:
	//
	//     myApp --flag value
	//
	// An Option can only have one long flag defined so any previously defined
	// value will be overwritten. This function does nothing if Parse has
	// already been called.
	LongFlag(name string)

	// EnvVar defines the name of an environment variable that can be used to
	// set this option. The string value of the environment variable must be
	// convertible to the type of the option.
	EnvVar(name string)

	// ConfigKey defines the key this option can use to set itself from a JSON
	// configuration file. The value stored under this key must be convertible
	// to the Option Type.
	ConfigKey(name string)

	// Default defines a value that will be used by the option if no flags,
	// environment variables, or configuration file values are set or found. The
	// value must be the same type as the Option.
	Default(value interface{})

	// Value returns the encapsulated value held by this Option. Value should
	// not be called until Parse has been called. In general Value only needs to
	// be called by a parent Options type.
	Value() value.Data

	// RegisterFlags causes the flags defined on this option to be registered
	// with the provided FlagSet. In general RegisterFlags only needs to be
	// called by a parent Options type.
	RegisterFlags(flags *flag.FlagSet) error

	// Parse this option with the provided config. In general Parse only needs
	// to be called by a parent Options type.
	Parse(config map[string]interface{}) error
}

type option struct {
	shortFlag   rune
	longFlag    string
	envVar      string
	configKey   string
	description string

	short  value.Data
	long   value.Data
	env    value.Data
	config value.Data

	flags  *flag.FlagSet
	typeOf reflect.Type

	backstop interface{}
	pointer  interface{}
}

// NewOption returns an option with i being a pointer to a variable of the type
// of this option (*bool, *int, *int64, *uint, *uint64, *float64, *string,
// *time.Duration). The description is used when producing usage information.
// Providing an invalid type for i will not error here, but will generate an
// error when the Option is parsed.
func NewOption(i interface{}, description string) Option {
	if i == nil {
		return &option{description:description}
	}

	typeOf := reflect.Indirect(reflect.ValueOf(i)).Type()
	return &option{
		description: description, typeOf: typeOf, pointer: i}
}

func (o *option) Flags(short rune, longFlag string) {
	o.shortFlag = short
	o.longFlag = longFlag
}

func (o *option) ShortFlag(name rune) {
	o.shortFlag = name
}

func (o *option) LongFlag(name string) {
	o.longFlag = name
}

func (o *option) EnvVar(name string) {
	o.envVar = name
}

func (o *option) ConfigKey(name string) {
	o.configKey = name
}

func (o *option) Default(value interface{}) {
	o.backstop = value
}

func (o *option) Value() value.Data {
	var v value.Data

	if o.flags != nil {
		o.flags.Visit(func(f *flag.Flag) {
			if !v.Set && f.Name == string(o.shortFlag) {
				v = o.short
			}

			if !v.Set && f.Name == o.longFlag {
				v = o.long
			}
		})
	}

	if !v.Set && o.env.Set {
		v = o.env
	}

	if !v.Set && o.config.Set {
		v = o.config
	}

	if !v.Set {
		v = value.New(o.backstop)
	}

	return v
}

func (o *option) RegisterFlags(flags *flag.FlagSet) error {
	o.flags = flags

	if o.flags == nil {
		return nil
	}

	if o.shortFlag != 0 {
		v, err := o.registerFlag(string(o.shortFlag))

		if err != nil {
			return fmt.Errorf("failed to set short flag: %s", err)
		}

		o.short = v
	}

	if o.longFlag != "" {
		v, err := o.registerFlag(o.longFlag)

		if err != nil {
			return fmt.Errorf("failed to set long flag: %s", err)
		}

		o.long = v
	}

	return nil
}

func (o *option) Parse(config map[string]interface{}) error {
	var err error

	if o.pointer == nil {
		return fmt.Errorf("option with description '%s' "+
			"not registered with a value", o.description)
	}

	if reflect.TypeOf(o.pointer).Kind() != reflect.Ptr {
		return fmt.Errorf("type Option requires pointer value(*%s, not %[1]s)",
			o.typeOf.String())
	}

	if err = o.setConfig(config); err != nil {
		return fmt.Errorf("failed to parse option config: %s", err)
	}

	if env := os.Getenv(o.envVar); env != "" {
		if o.env, err = value.Coerce(env, o.pointer); err != nil {
			return fmt.Errorf("failed to parse environment option '%s': %s",
				o.envVar, err)
		}
	}

	if err = o.Value().AssignTo(o.pointer); err != nil {
		return fmt.Errorf("failed to set option: %s", err)
	}

	return nil
}

func (o *option) registerFlag(name string) (value.Data, error) {
	var f func() interface{}
	var ok bool

	switch o.pointer.(type) {
	case *bool:
		v, success := o.backstop.(bool)
		f = func() interface{} { return o.flags.Bool(name, v, o.description) }
		ok = success
	case *int:
		v, success := o.backstop.(int)
		f = func() interface{} { return o.flags.Int(name, v, o.description) }
		ok = success
	case *int64:
		v, success := o.backstop.(int64)
		f = func() interface{} { return o.flags.Int64(name, v, o.description) }
		ok = success
	case *uint:
		v, success := o.backstop.(uint)
		f = func() interface{} { return o.flags.Uint(name, v, o.description) }
		ok = success
	case *uint64:
		v, success := o.backstop.(uint64)
		f = func() interface{} { return o.flags.Uint64(name, v, o.description) }
		ok = success
	case *float64:
		v, success := o.backstop.(float64)
		f = func() interface{} { return o.flags.Float64(name, v, o.description) }
		ok = success
	case *string:
		v, success := o.backstop.(string)
		f = func() interface{} { return o.flags.String(name, v, o.description) }
		ok = success
	case *time.Duration:
		v, success := o.backstop.(time.Duration)
		f = func() interface{} { return o.flags.Duration(name, v, o.description) }
		ok = success
	default:
		return value.Data{}, fmt.Errorf("invalid option type: %T", o.pointer)
	}

	if !ok && o.backstop != nil {
		return value.Data{}, fmt.Errorf(
			"cannot use default option %v (%[1]T) as %T for flag %s",
			o.backstop, o.pointer, name)
	}

	return value.New(f()), nil

}

func (o *option) setConfig(config map[string]interface{}) error {
	v, ok := config[o.configKey]

	if ok && !reflect.TypeOf(v).ConvertibleTo(o.typeOf) {
		return fmt.Errorf("cannot convert config type %T to %s for '%s'",
			v, o.typeOf.String(), o.configKey)
	} else if ok {
		o.config = value.New(v)
	}

	return nil
}
