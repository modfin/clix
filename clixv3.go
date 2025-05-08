package clix

import (
	"time"
)

// CommandReaderV3 is the interface
// that a v3 cli.Command must support in order to be convertible to a v2 ContextReader and the use of clix.Parse[]()
type CommandReaderV3 interface {
	String(name string) string
	Int(name string) int
	Uint(name string) uint
	Bool(name string) bool
	Float(name string) float64
	Timestamp(name string) time.Time
	Duration(name string) time.Duration
	StringSlice(name string) []string
	IntSlice(name string) []int
	UintSlice(name string) []uint
	FloatSlice(name string) []float64
}

// V3 converts a v3 CommandReaderV3 to a v2 ContextReader
// example
//
//	 func(ctx context.Context, cmd *cli.Command) error {
//		  config := clix.Parse[Config](clix.V3(cmd))
func V3(cmd CommandReaderV3) ContextReader {
	return &proxy3to2{c: cmd}
}

// ParseCommand converts a v3 CommandReaderV3 to a v2 ContextReader and uses the default Parse command,
// it is just a shorthand for `clix.Parse[Config](clix.V3(cmd))`
// example
//
//	 func(ctx context.Context, cmd *cli.Command) error {
//		  config := clix.ParseCommand[Config](cmd)
func ParseCommand[A any](cmd CommandReaderV3) A {
	return Parse[A](V3(cmd))
}

type proxy3to2 struct {
	c CommandReaderV3
}

func (p proxy3to2) String(name string) string {
	return p.c.String(name)
}

func (p proxy3to2) Int(name string) int {
	return p.c.Int(name)
}

func (p proxy3to2) Int64(name string) int64 {
	return int64(p.c.Int(name))
}

func (p proxy3to2) Uint(name string) uint {
	return p.c.Uint(name)
}

func (p proxy3to2) Uint64(name string) uint64 {
	return uint64(p.c.Uint(name))

}

func (p proxy3to2) Bool(name string) bool {
	return p.c.Bool(name)
}

func (p proxy3to2) Float64(name string) float64 {
	return p.c.Float(name)
}

func (p proxy3to2) Timestamp(name string) *time.Time {
	t := p.c.Timestamp(name)
	if t.IsZero() {
		return nil
	}
	return &t
}

func (p proxy3to2) Duration(name string) time.Duration {
	return p.c.Duration(name)
}

func (p proxy3to2) StringSlice(name string) []string {
	return p.c.StringSlice(name)
}

func (p proxy3to2) IntSlice(name string) []int {
	return p.c.IntSlice(name)
}

func (p proxy3to2) Int64Slice(name string) []int64 {
	return toInt64Slice(p.c.IntSlice(name))
}

func (p proxy3to2) UintSlice(name string) []uint {
	return p.c.UintSlice(name)
}

func (p proxy3to2) Uint64Slice(name string) []uint64 {
	return toUint64Slice(p.c.UintSlice(name))
}

func (p proxy3to2) Float64Slice(name string) []float64 {
	return p.c.FloatSlice(name)
}

func toInt64Slice(sl []int) []int64 {
	ints := make([]int64, len(sl))
	for i, v := range sl {
		ints[i] = int64(v)
	}
	return ints
}

func toUint64Slice(sl []uint) []uint64 {
	ints := make([]uint64, len(sl))
	for i, v := range sl {
		ints[i] = uint64(v)
	}
	return ints
}
