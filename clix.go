package clix

import (
	"reflect"
	"time"
)

type Context interface {
	Timestamp(name string) *time.Time
	Duration(name string) time.Duration
	String(name string) string
	Int(name string) int
	Int64(name string) int64
	Uint(name string) uint
	Uint64(name string) uint64
	Bool(name string) bool
	Float64(name string) float64

	StringSlice(name string) []string
	IntSlice(name string) []int
	Int64Slice(name string) []int64
	UintSlice(name string) []uint
	Uint64Slice(name string) []uint64
	Float64Slice(name string) []float64
}

// Parse takes a cli.Context or struct with same interface and a config struct type.
// It fills in the fields of the config struct with values from the Context. The config
// struct fields must have a "cli" tag with the name of the flag to take the value from.
// If the field is a struct, the function will recursively fill in the fields of the
// struct. If the field is a slice, the function will fill in the slice with
// the values of the flag. The function will return the config struct with all
// the fields filled in.
func Parse[A any](c Context) A {

	var cfg A

	var AssignValueToCliFields func(v interface{})
	AssignValueToCliFields = func(v interface{}) {
		// Get the reflection value of the input struct
		val := reflect.ValueOf(v).Elem()

		// Iterate over the struct fields
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := val.Type().Field(i)

			// Check if the field is a struct

			tag := fieldType.Tag.Get("cli")

			if tag == "" && field.Kind() == reflect.Struct {
				if field.Addr().CanInterface() { // is it a public struct?
					AssignValueToCliFields(field.Addr().Interface())
				}
				continue
			}

			if tag != "" {

				if field.Type() == reflect.TypeOf(time.Time{}) || field.Type() == reflect.PointerTo(reflect.TypeOf(time.Time{})) {
					t := c.Timestamp(tag)
					if t != nil {
						if field.Kind() == reflect.Ptr {
							field.Set(reflect.ValueOf(t))
						} else {
							field.Set(reflect.ValueOf(*t))
						}
					}
					continue
				}

				if field.Type() == reflect.TypeOf(time.Duration(0)) {
					field.Set(reflect.ValueOf(c.Duration(tag)))
					continue
				}

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
			}
		}
	}
	AssignValueToCliFields(&cfg)

	return cfg
}
