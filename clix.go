package clix

import (
	"reflect"
	"time"
)

// ContextReader defines the interface for reading values from a CLI context.
// This abstraction allows for easier testing by allowing mock implementations.
type ContextReader interface {
	String(name string) string
	Int(name string) int
	Int64(name string) int64
	Uint(name string) uint
	Uint64(name string) uint64
	Bool(name string) bool
	Float64(name string) float64
	Timestamp(name string) *time.Time
	Duration(name string) time.Duration
	StringSlice(name string) []string
	IntSlice(name string) []int
	Int64Slice(name string) []int64
	UintSlice(name string) []uint
	Uint64Slice(name string) []uint64
	Float64Slice(name string) []float64
}

// Parse converts CLI context into a typed configuration struct.
// It uses reflection to map CLI flags to struct fields based on struct tags.
// Usage:
//
//	type Config struct {
//	    Host     string        `cli:"host"`
//	    Port     int           `cli:"port"`
//	    Timeout  time.Duration `cli:"timeout"`
//	    Database struct {
//	        Name string `cli:"name"`
//	        User string `cli:"user"`
//	    } `cli-prefix:"db-"`
//	}
//
//	cfg := clix.Parse[Config](ctx)
func Parse[A any](c ContextReader) A {
	var cfg A
	AssignValueToCliFields(&cfg, "", c)
	return cfg
}

// AssignValueToCliFields recursively assigns values from CLI flags to struct fields.
// It uses reflection to iterate over the struct fields and set their values based on CLI flags.
// Parameters:
//   - v: a pointer to the struct to populate
//   - prefix: prefix for the CLI flags (used for nested structs)
//   - c: the CLI context containing the flag values
func AssignValueToCliFields(v interface{}, prefix string, c ContextReader) {
	// Get the reflection value of the input struct
	val := reflect.ValueOf(v).Elem()

	// Iterate over the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// Skip unexported fields (they start with lowercase)
		if !field.CanSet() {
			continue
		}

		// Get the "cli" tag value
		tag := fieldType.Tag.Get("cli")

		// Handle nested structs without a cli tag
		if tag == "" && field.Kind() == reflect.Struct {
			// Get the prefix for the nested struct
			nestedPrefix := fieldType.Tag.Get("cli-prefix")
			// Process the nested struct recursively if it's addressable
			if field.Addr().CanInterface() {
				AssignValueToCliFields(field.Addr().Interface(), prefix+nestedPrefix, c)
			}
			continue
		}

		// Process fields with cli tags
		if tag != "" {
			// Combine the prefix with the tag
			fullTag := prefix + tag

			// Handle time.Time and *time.Time types
			if field.Type() == reflect.TypeOf(time.Time{}) ||
				field.Type() == reflect.PointerTo(reflect.TypeOf(time.Time{})) {
				setTimeValue(c, fullTag, field)
				continue
			}

			// Handle time.Duration type
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				field.Set(reflect.ValueOf(c.Duration(fullTag)))
				continue
			}

			// Handle other types based on their Kind
			setFieldValue(c, fullTag, field)
		}
	}
}

// setTimeValue handles setting time.Time values from CLI flags.
// It handles both time.Time and *time.Time types.
func setTimeValue(c ContextReader, tag string, field reflect.Value) {
	t := c.Timestamp(tag)
	if t != nil {
		if field.Kind() == reflect.Ptr {
			field.Set(reflect.ValueOf(t))
		} else {
			field.Set(reflect.ValueOf(*t))
		}
	}
}

// setFieldValue sets the value of a field based on its Kind.
// It handles primitive types and slices of primitive types.
func setFieldValue(c ContextReader, tag string, field reflect.Value) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(c.String(tag))
	case reflect.Int:
		field.SetInt(int64(c.Int(tag)))
	case reflect.Int64:
		field.SetInt(c.Int64(tag))
	case reflect.Uint:
		field.SetUint(uint64(c.Uint(tag)))
	case reflect.Uint64:
		field.SetUint(c.Uint64(tag))
	case reflect.Bool:
		field.SetBool(c.Bool(tag))
	case reflect.Float64:
		field.SetFloat(c.Float64(tag))
	case reflect.Slice:
		setSliceValue(c, tag, field)
	}
}

// setSliceValue handles setting slice values from CLI flags.
// It supports various slice types like []string, []int, etc.
func setSliceValue(c ContextReader, tag string, field reflect.Value) {
	switch field.Type() {
	case reflect.TypeOf([]string{}):
		field.Set(reflect.ValueOf(c.StringSlice(tag)))
	case reflect.TypeOf([]int{}):
		field.Set(reflect.ValueOf(c.IntSlice(tag)))
	case reflect.TypeOf([]int64{}):
		field.Set(reflect.ValueOf(c.Int64Slice(tag)))
	case reflect.TypeOf([]uint{}):
		field.Set(reflect.ValueOf(c.UintSlice(tag)))
	case reflect.TypeOf([]uint64{}):
		field.Set(reflect.ValueOf(c.Uint64Slice(tag)))
	case reflect.TypeOf([]float64{}):
		field.Set(reflect.ValueOf(c.Float64Slice(tag)))
	}
}
