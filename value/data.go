package value

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Data holds an untyped (interface{}) value which can be assigned to a bool,
// int, int64, uint, uint64, float64, string, or time.Duration.
type Data struct {
	Set     bool
	pointer interface{}
}

// New creates a new Data type with the given value.
func New(data interface{}) Data {
	return Data{Set: true, pointer: data}
}

// Coerce the given string into a Data type holding a value of typeOf.
func Coerce(data string, typeOf interface{}) (Data, error) {
	var r interface{}
	var err error

	switch typeOf.(type) {
	case *bool:
		r, err = strconv.ParseBool(data)
	case *int:
		v, e := strconv.ParseInt(data, 10, 64)
		r, err = int(v), e
	case *int64:
		r, err = strconv.ParseInt(data, 10, 64)
	case *uint:
		v, e := strconv.ParseUint(data, 10, 64)
		r, err = uint(v), e
	case *uint64:
		r, err = strconv.ParseUint(data, 10, 64)
	case *float64:
		r, err = strconv.ParseFloat(data, 64)
	case *string:
		r = data
	case *time.Duration:
		r, err = time.ParseDuration(data)
	default:
		err = fmt.Errorf("invalid type: %T", typeOf)
	}

	if err != nil {
		err = fmt.Errorf("value.Data: cannot coerce '%s': %s", data, err)
	}

	return New(r), err
}

// Pointer returns the underlying value wrapped by this data type.
func (d Data) Pointer() interface{} {
	return d.pointer
}

// AssignTo sets the given pointer to point to the value held by this Data type.
func (d Data) AssignTo(pointer interface{}) (err error) {
	var to, from reflect.Value
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("value.Data %v '%v', failed to assign to type %v",
				from.Type(), from.Interface(), to.Type())
		}
	}()

	if pointer == nil || d.pointer == nil {
		return nil
	}

	from = reflect.ValueOf(d.pointer)
	to = reflect.ValueOf(pointer)

	if from.Kind() == reflect.Ptr {
		from = reflect.Indirect(from)
	}

	if to.Kind() == reflect.Ptr {
		to = reflect.Indirect(to)
	} else {
		return fmt.Errorf("cannot assign to %v, should be *%[1]v", to.Type())
	}

	data := from.Convert(to.Type())

	switch p := pointer.(type) {
	case *bool:
		*p = data.Bool()
	case *int:
		*p = int(data.Int())
	case *int64:
		*p = data.Int()
	case *uint:
		*p = uint(data.Uint())
	case *uint64:
		*p = data.Uint()
	case *float64:
		*p = data.Float()
	case *string:
		*p = data.String()
	case *time.Duration:
		*p = time.Duration(data.Int())
	default:
		return fmt.Errorf("value.Data invalid pointer type: %T", p)
	}

	return nil
}
