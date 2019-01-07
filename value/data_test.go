package value_test

import (
	"fmt"
	"github.com/domdavis/goconfigure/value"
	"strings"
	"testing"
	"time"
)

func ExampleCoerce() {
	var typeOf string
	if v, err := value.Coerce("text", &typeOf); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v)
	}

	// Output:
	// {true text}
}

func ExampleData_AssignTo() {
	var v string
	d := value.New("example")
	if err := d.AssignTo(&v); err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)

	// Output:
	// example
}

func TestCoerce(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		var p bool
		if v, err := value.Coerce("t", &p); err != nil {
			t.Errorf("Unexpected error coercing bool: %s", err.Error())
		} else if v.Pointer() != true {
			t.Errorf("Unexpected value coercing bool: %v", v.Pointer())
		}
	})

	t.Run("int", func(t *testing.T) {
		var p int
		if v, err := value.Coerce("1", &p); err != nil {
			t.Errorf("Unexpected error coercing int: %s", err.Error())
		} else if v.Pointer() != 1 {
			t.Errorf("Unexpected value coercing int: %v", v.Pointer())
		}
	})

	t.Run("int64", func(t *testing.T) {
		var p int64
		if v, err := value.Coerce("1", &p); err != nil {
			t.Errorf("Unexpected error coercing int64: %s", err.Error())
		} else if v.Pointer() != int64(1) {
			t.Errorf("Unexpected value coercing int64: %v", v.Pointer())
		}
	})

	t.Run("uint", func(t *testing.T) {
		var p uint
		if v, err := value.Coerce("1", &p); err != nil {
			t.Errorf("Unexpected error coercing uint: %s", err.Error())
		} else if v.Pointer() != uint(1) {
			t.Errorf("Unexpected value coercing uint: %v", v.Pointer())
		}
	})

	t.Run("unit64", func(t *testing.T) {
		var p uint64
		if v, err := value.Coerce("1", &p); err != nil {
			t.Errorf("Unexpected error coercing uint64: %s", err.Error())
		} else if v.Pointer() != uint64(1) {
			t.Errorf("Unexpected value coercing uint64: %v", v.Pointer())
		}
	})

	t.Run("float64", func(t *testing.T) {
		var p float64
		if v, err := value.Coerce("1", &p); err != nil {
			t.Errorf("Unexpected error coercing float64: %s", err.Error())
		} else if v.Pointer() != float64(1) {
			t.Errorf("Unexpected value coercing float64: %v", v.Pointer())
		}
	})

	t.Run("string", func(t *testing.T) {
		var p string
		if v, err := value.Coerce("text", &p); err != nil {
			t.Errorf("Unexpected error coercing string: %s", err.Error())
		} else if v.Pointer() != "text" {
			t.Errorf("Unexpected value coercing string: %v", v.Pointer())
		}
	})

	t.Run("time.Duration", func(t *testing.T) {
		var p time.Duration
		if v, err := value.Coerce("1ns", &p); err != nil {
			t.Errorf("Unexpected error coercing time.Duration: %s", err.Error())
		} else if v.Pointer() != time.Nanosecond {
			t.Errorf("Unexpected value coercing time.Duration: %v", v.Pointer())
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var p []string
		_, err := value.Coerce("invalid", &p)

		expected := "value.Data: cannot coerce 'invalid': " +
			"invalid type: *[]string"
		if err.Error() != expected {
			t.Errorf("Unexpected error coercing invalid value: %s", err.Error())
		}
	})

	const prefix = "value.Data: cannot coerce 'invalid':"
	for _, p := range []interface{} {
		new(bool),
		new(int),
		new(int64),
		new(uint),
		new(uint64),
		new(float64),
		"",
		time.Duration(0),
		[]string{},
	} {
		test := fmt.Sprintf("invalid %T", p)
		t.Run(test, func(t *testing.T) {
			_, err := value.Coerce("invalid", &p)
			if err == nil || !strings.HasPrefix(err.Error(), prefix) {
				t.Errorf("Unexpected error coercing invalid %T: %s",
					p, err.Error())
			}
		})
	}
}

func TestData_AssignTo(t *testing.T) {
	t.Run("nil data will not error", func(t *testing.T) {
		var p interface{}
		d := value.New(nil)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning from nil value: %s",
				err.Error())
		}
	})

	t.Run("nil assignable pointer will not error", func(t *testing.T) {
		d := value.New("test")
		if err := d.AssignTo(nil); err != nil {
			t.Errorf("Unexpected error assigning to nil value: %s", err.Error())
		}
	})

	t.Run("assignable must be pointer", func(t *testing.T) {
		var p string
		d := value.New("test")
		err := d.AssignTo(p)

		expected := "cannot assign to string, should be *string"
		if err == nil || err.Error() != expected {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		}
	})

	t.Run("assignTo works with pointers", func(t *testing.T) {
		var p string
		v := "test"
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from pointer", v, p)
		}
	})

	t.Run("assignTo works with bool", func(t *testing.T) {
		var p bool
		v := true
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from bool", v, p)
		}
	})

	t.Run("assignTo works with int", func(t *testing.T) {
		var p int
		v := 1
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from int", v, p)
		}
	})

	t.Run("assignTo works with int64", func(t *testing.T) {
		var p int64
		v := int64(1)
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from int64", v, p)
		}
	})

	t.Run("assignTo works with uint", func(t *testing.T) {
		var p uint
		v := uint(1)
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from uint", v, p)
		}
	})

	t.Run("assignTo works with uint64", func(t *testing.T) {
		var p uint64
		v := uint64(1)
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from uint64", v, p)
		}
	})

	t.Run("assignTo works with float64", func(t *testing.T) {
		var p float64
		v := float64(1)
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from float64", v, p)
		}
	})

	t.Run("assignTo works with string", func(t *testing.T) {
		var p string
		v := "test"
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from string", v, p)
		}
	})

	t.Run("assignTo works with Duration", func(t *testing.T) {
		var p time.Duration
		v := time.Duration(1)
		d := value.New(&v)
		if err := d.AssignTo(&p); err != nil {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		} else if p != v {
			t.Errorf("Expected '%v', got '%v' assigning from Duration", v, p)
		}
	})

	t.Run("assignTo fails assign to an invalid type", func(t *testing.T) {
		var p []string
		v := true
		d := value.New(&v)
		err := d.AssignTo(&p)

		expected := "value.Data bool 'true', failed to assign to type []string"
		if err == nil || err.Error() != expected {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		}
	})

	t.Run("assignTo fails assign from an invalid type", func(t *testing.T) {
		var p []string
		var v []string
		d := value.New(&v)
		err := d.AssignTo(&p)

		expected := "value.Data invalid pointer type: *[]string"
		if err == nil || err.Error() != expected {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		}
	})

	t.Run("assignTo will not cast non string values", func(t *testing.T) {
		var p string
		var v int
		d := value.New(&v)
		err := d.AssignTo(&p)

		expected := "value.Data: invalid cast of int to string"
		if err == nil || err.Error() != expected {
			t.Errorf("Unexpected error assigning to value: %s", err.Error())
		}
	})
}

