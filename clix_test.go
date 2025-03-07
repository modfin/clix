package clix

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test data structures
type BasicConfig struct {
	StringVal  string  `cli:"string-val"`
	IntVal     int     `cli:"int-val"`
	Int64Val   int64   `cli:"int64-val"`
	UintVal    uint    `cli:"uint-val"`
	Uint64Val  uint64  `cli:"uint64-val"`
	BoolVal    bool    `cli:"bool-val"`
	FloatVal   float64 `cli:"float-val"`
	unexported string  `cli:"unexported"` // Should be skipped
}

type TimeConfig struct {
	TimeVal     time.Time     `cli:"time-val"`
	TimePointer *time.Time    `cli:"time-ptr"`
	DurationVal time.Duration `cli:"duration-val"`
	unexported  time.Time     `cli:"unexported-time"` // Should be skipped
}

type SliceConfig struct {
	StringSlice  []string  `cli:"string-slice"`
	IntSlice     []int     `cli:"int-slice"`
	Int64Slice   []int64   `cli:"int64-slice"`
	UintSlice    []uint    `cli:"uint-slice"`
	Uint64Slice  []uint64  `cli:"uint64-slice"`
	Float64Slice []float64 `cli:"float64-slice"`
}

type NestedConfig struct {
	TopLevel string `cli:"top-level"`
	Database struct {
		Host string `cli:"host"`
		Port int    `cli:"port"`
	} `cli-prefix:"db-"`
	Advanced struct {
		Feature1 bool `cli:"feature1"`
		Feature2 bool `cli:"feature2"`
	} `cli-prefix:"adv-"`
}

type MixedConfig struct {
	Name     string        `cli:"name"`
	Timeout  time.Duration `cli:"timeout"`
	Enabled  bool          `cli:"enabled"`
	Database struct {
		Host     string   `cli:"host"`
		Port     int      `cli:"port"`
		Replicas []string `cli:"replicas"`
	} `cli-prefix:"db-"`
}

// Custom mock that implements the minimum interface needed for testing
type cliContextMock struct {
	stringMap       map[string]string
	intMap          map[string]int
	int64Map        map[string]int64
	uintMap         map[string]uint
	uint64Map       map[string]uint64
	boolMap         map[string]bool
	float64Map      map[string]float64
	timestampMap    map[string]*time.Time
	durationMap     map[string]time.Duration
	stringSliceMap  map[string][]string
	intSliceMap     map[string][]int
	int64SliceMap   map[string][]int64
	uintSliceMap    map[string][]uint
	uint64SliceMap  map[string][]uint64
	float64SliceMap map[string][]float64
}

func newMockContext() *cliContextMock {
	return &cliContextMock{
		stringMap:       make(map[string]string),
		intMap:          make(map[string]int),
		int64Map:        make(map[string]int64),
		uintMap:         make(map[string]uint),
		uint64Map:       make(map[string]uint64),
		boolMap:         make(map[string]bool),
		float64Map:      make(map[string]float64),
		timestampMap:    make(map[string]*time.Time),
		durationMap:     make(map[string]time.Duration),
		stringSliceMap:  make(map[string][]string),
		intSliceMap:     make(map[string][]int),
		int64SliceMap:   make(map[string][]int64),
		uintSliceMap:    make(map[string][]uint),
		uint64SliceMap:  make(map[string][]uint64),
		float64SliceMap: make(map[string][]float64),
	}
}

// Implementation of cli.Context methods used by our package
func (m *cliContextMock) String(name string) string          { return m.stringMap[name] }
func (m *cliContextMock) Int(name string) int                { return m.intMap[name] }
func (m *cliContextMock) Int64(name string) int64            { return m.int64Map[name] }
func (m *cliContextMock) Uint(name string) uint              { return m.uintMap[name] }
func (m *cliContextMock) Uint64(name string) uint64          { return m.uint64Map[name] }
func (m *cliContextMock) Bool(name string) bool              { return m.boolMap[name] }
func (m *cliContextMock) Float64(name string) float64        { return m.float64Map[name] }
func (m *cliContextMock) Timestamp(name string) *time.Time   { return m.timestampMap[name] }
func (m *cliContextMock) Duration(name string) time.Duration { return m.durationMap[name] }
func (m *cliContextMock) StringSlice(name string) []string   { return m.stringSliceMap[name] }
func (m *cliContextMock) IntSlice(name string) []int         { return m.intSliceMap[name] }
func (m *cliContextMock) Int64Slice(name string) []int64     { return m.int64SliceMap[name] }
func (m *cliContextMock) UintSlice(name string) []uint       { return m.uintSliceMap[name] }
func (m *cliContextMock) Uint64Slice(name string) []uint64   { return m.uint64SliceMap[name] }
func (m *cliContextMock) Float64Slice(name string) []float64 { return m.float64SliceMap[name] }

func TestParseBasicTypes(t *testing.T) {
	// Create a mock context with basic flag values
	ctx := newMockContext()
	ctx.stringMap["string-val"] = "hello"
	ctx.intMap["int-val"] = 42
	ctx.int64Map["int64-val"] = 9223372036854775807 // Max int64
	ctx.uintMap["uint-val"] = 123
	ctx.uint64Map["uint64-val"] = 18446744073709551615 // Max uint64
	ctx.boolMap["bool-val"] = true
	ctx.float64Map["float-val"] = 3.14159
	ctx.stringMap["unexported"] = "should-be-ignored"

	// Parse the context into a BasicConfig
	config := Parse[BasicConfig](ctx)

	// Verify the values were correctly parsed
	assert.Equal(t, "hello", config.StringVal)
	assert.Equal(t, 42, config.IntVal)
	assert.Equal(t, int64(9223372036854775807), config.Int64Val)
	assert.Equal(t, uint(123), config.UintVal)
	assert.Equal(t, uint64(18446744073709551615), config.Uint64Val)
	assert.Equal(t, true, config.BoolVal)
	assert.Equal(t, 3.14159, config.FloatVal)
	assert.Equal(t, "", config.unexported) // Unexported field should be untouched
}

func TestParseTimeTypes(t *testing.T) {
	// Sample time for testing
	sampleTime, _ := time.Parse(time.RFC3339, "2023-01-02T15:04:05Z")

	// Create a mock context with time-related flag values
	ctx := newMockContext()
	ctx.timestampMap["time-val"] = &sampleTime
	ctx.timestampMap["time-ptr"] = &sampleTime
	ctx.durationMap["duration-val"] = 90 * time.Minute // 1h30m

	// Create another time for the unexported field (which should be ignored)
	unexportedTime, _ := time.Parse(time.RFC3339, "2020-12-31T23:59:59Z")
	ctx.timestampMap["unexported-time"] = &unexportedTime

	// Parse the context into a TimeConfig
	config := Parse[TimeConfig](ctx)

	// Verify the values were correctly parsed
	assert.Equal(t, sampleTime, config.TimeVal)
	assert.NotNil(t, config.TimePointer)
	assert.Equal(t, sampleTime, *config.TimePointer)
	assert.Equal(t, 90*time.Minute, config.DurationVal)
	assert.Equal(t, time.Time{}, config.unexported) // Unexported field should be untouched
}

func TestParseSliceTypes(t *testing.T) {
	// Create a mock context with slice flag values
	ctx := newMockContext()
	ctx.stringSliceMap["string-slice"] = []string{"a", "b", "c"}
	ctx.intSliceMap["int-slice"] = []int{1, 2, 3}
	ctx.int64SliceMap["int64-slice"] = []int64{4, 5, 6}
	ctx.uintSliceMap["uint-slice"] = []uint{7, 8, 9}
	ctx.uint64SliceMap["uint64-slice"] = []uint64{10, 11, 12}
	ctx.float64SliceMap["float64-slice"] = []float64{1.1, 2.2, 3.3}

	// Parse the context into a SliceConfig
	config := Parse[SliceConfig](ctx)

	// Verify the values were correctly parsed
	assert.Equal(t, []string{"a", "b", "c"}, config.StringSlice)
	assert.Equal(t, []int{1, 2, 3}, config.IntSlice)
	assert.Equal(t, []int64{4, 5, 6}, config.Int64Slice)
	assert.Equal(t, []uint{7, 8, 9}, config.UintSlice)
	assert.Equal(t, []uint64{10, 11, 12}, config.Uint64Slice)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, config.Float64Slice)
}

func TestParseNestedStructs(t *testing.T) {
	// Create a mock context with nested flag values
	ctx := newMockContext()
	ctx.stringMap["top-level"] = "root"
	ctx.stringMap["db-host"] = "localhost"
	ctx.intMap["db-port"] = 5432
	ctx.boolMap["adv-feature1"] = true
	ctx.boolMap["adv-feature2"] = false

	// Parse the context into a NestedConfig
	config := Parse[NestedConfig](ctx)

	// Verify the values were correctly parsed
	assert.Equal(t, "root", config.TopLevel)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, true, config.Advanced.Feature1)
	assert.Equal(t, false, config.Advanced.Feature2)
}

func TestParseMixedConfig(t *testing.T) {
	// Create a mock context with mixed flag values
	ctx := newMockContext()
	ctx.stringMap["name"] = "my-app"
	ctx.durationMap["timeout"] = 2*time.Minute + 30*time.Second
	ctx.boolMap["enabled"] = true
	ctx.stringMap["db-host"] = "db.example.com"
	ctx.intMap["db-port"] = 3306
	ctx.stringSliceMap["db-replicas"] = []string{"replica1", "replica2", "replica3"}

	// Parse the context into a MixedConfig
	config := Parse[MixedConfig](ctx)

	// Verify the values were correctly parsed
	assert.Equal(t, "my-app", config.Name)
	assert.Equal(t, 2*time.Minute+30*time.Second, config.Timeout)
	assert.Equal(t, true, config.Enabled)
	assert.Equal(t, "db.example.com", config.Database.Host)
	assert.Equal(t, 3306, config.Database.Port)
	assert.Equal(t, []string{"replica1", "replica2", "replica3"}, config.Database.Replicas)
}

func TestMissingValues(t *testing.T) {
	// Create an empty mock context
	ctx := newMockContext()

	// Parse the context into configs
	basicConfig := Parse[BasicConfig](ctx)
	timeConfig := Parse[TimeConfig](ctx)
	sliceConfig := Parse[SliceConfig](ctx)

	// Verify default values are used
	assert.Equal(t, "", basicConfig.StringVal)
	assert.Equal(t, 0, basicConfig.IntVal)
	assert.Equal(t, int64(0), basicConfig.Int64Val)
	assert.Equal(t, uint(0), basicConfig.UintVal)
	assert.Equal(t, uint64(0), basicConfig.Uint64Val)
	assert.Equal(t, false, basicConfig.BoolVal)
	assert.Equal(t, 0.0, basicConfig.FloatVal)

	assert.Equal(t, time.Time{}, timeConfig.TimeVal)
	assert.Nil(t, timeConfig.TimePointer)
	assert.Equal(t, time.Duration(0), timeConfig.DurationVal)

	assert.Nil(t, sliceConfig.StringSlice)
	assert.Nil(t, sliceConfig.IntSlice)
	assert.Nil(t, sliceConfig.Int64Slice)
	assert.Nil(t, sliceConfig.UintSlice)
	assert.Nil(t, sliceConfig.Uint64Slice)
	assert.Nil(t, sliceConfig.Float64Slice)
}
