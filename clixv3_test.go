package clix

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockCommandReaderV3 implements the CommandReaderV3 interface for testing
type mockCommandReaderV3 struct {
	stringMap      map[string]string
	intMap         map[string]int
	uintMap        map[string]uint
	boolMap        map[string]bool
	floatMap       map[string]float64
	timestampMap   map[string]time.Time
	durationMap    map[string]time.Duration
	stringSliceMap map[string][]string
	intSliceMap    map[string][]int64
	uintSliceMap   map[string][]uint64
	floatSliceMap  map[string][]float64
}

func (m *mockCommandReaderV3) String(name string) string          { return m.stringMap[name] }
func (m *mockCommandReaderV3) Int(name string) int                { return m.intMap[name] }
func (m *mockCommandReaderV3) Uint(name string) uint              { return m.uintMap[name] }
func (m *mockCommandReaderV3) Bool(name string) bool              { return m.boolMap[name] }
func (m *mockCommandReaderV3) Float(name string) float64          { return m.floatMap[name] }
func (m *mockCommandReaderV3) Timestamp(name string) time.Time    { return m.timestampMap[name] }
func (m *mockCommandReaderV3) Duration(name string) time.Duration { return m.durationMap[name] }
func (m *mockCommandReaderV3) StringSlice(name string) []string   { return m.stringSliceMap[name] }
func (m *mockCommandReaderV3) IntSlice(name string) []int64       { return m.intSliceMap[name] }
func (m *mockCommandReaderV3) UintSlice(name string) []uint64     { return m.uintSliceMap[name] }
func (m *mockCommandReaderV3) FloatSlice(name string) []float64   { return m.floatSliceMap[name] }

// TestV3Conversion tests the V3 function that converts CommandReaderV3 to ContextReader
func TestV3Conversion(t *testing.T) {
	// Create a mock CommandReaderV3
	now := time.Now()
	cmdReader := &mockCommandReaderV3{
		stringMap:      map[string]string{"name": "test-name"},
		intMap:         map[string]int{"age": 30},
		uintMap:        map[string]uint{"count": 100},
		boolMap:        map[string]bool{"enabled": true},
		floatMap:       map[string]float64{"price": 99.99},
		timestampMap:   map[string]time.Time{"created": now},
		durationMap:    map[string]time.Duration{"timeout": 5 * time.Minute},
		stringSliceMap: map[string][]string{"tags": {"tag1", "tag2", "tag3"}},
		intSliceMap:    map[string][]int64{"values": {1, 2, 3}},
		uintSliceMap:   map[string][]uint64{"ids": {101, 102, 103}},
		floatSliceMap:  map[string][]float64{"rates": {1.1, 2.2, 3.3}},
	}

	// Convert the CommandReaderV3 to a ContextReader
	ctxReader := V3(cmdReader)

	// Test that the conversion works correctly for all types
	assert.Equal(t, "test-name", ctxReader.String("name"))
	assert.Equal(t, int(30), ctxReader.Int("age"))
	assert.Equal(t, int64(30), ctxReader.Int64("age"))
	assert.Equal(t, uint(100), ctxReader.Uint("count"))
	assert.Equal(t, uint64(100), ctxReader.Uint64("count"))
	assert.Equal(t, true, ctxReader.Bool("enabled"))
	assert.Equal(t, 99.99, ctxReader.Float64("price"))

	// For Timestamp, check if the conversion works correctly
	timestampPtr := ctxReader.Timestamp("created")
	assert.NotNil(t, timestampPtr)
	assert.Equal(t, now, *timestampPtr)

	assert.Equal(t, 5*time.Minute, ctxReader.Duration("timeout"))
	assert.Equal(t, []string{"tag1", "tag2", "tag3"}, ctxReader.StringSlice("tags"))

	// For slices, check if the conversion works correctly
	intSlice := ctxReader.IntSlice("values")
	assert.Equal(t, []int{1, 2, 3}, intSlice)

	int64Slice := ctxReader.Int64Slice("values")
	assert.Equal(t, []int64{1, 2, 3}, int64Slice)

	uintSlice := ctxReader.UintSlice("ids")
	assert.Equal(t, []uint{101, 102, 103}, uintSlice)

	uint64Slice := ctxReader.Uint64Slice("ids")
	assert.Equal(t, []uint64{101, 102, 103}, uint64Slice)

	floatSlice := ctxReader.Float64Slice("rates")
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, floatSlice)
}

// TestParseWithCommandReaderV3 tests using the Parse function with a CommandReaderV3
func TestParseWithCommandReaderV3(t *testing.T) {
	type Config struct {
		Name    string        `cli:"name"`
		Age     int           `cli:"age"`
		Enabled bool          `cli:"enabled"`
		Price   float64       `cli:"price"`
		Timeout time.Duration `cli:"timeout"`
		Tags    []string      `cli:"tags"`
		Values  []int         `cli:"values"`
		Rate    []float64     `cli:"rates"`
	}

	// Create a mock CommandReaderV3
	cmdReader := &mockCommandReaderV3{
		stringMap:      map[string]string{"name": "test-config"},
		intMap:         map[string]int{"age": 25},
		boolMap:        map[string]bool{"enabled": true},
		floatMap:       map[string]float64{"price": 49.99},
		durationMap:    map[string]time.Duration{"timeout": 10 * time.Minute},
		stringSliceMap: map[string][]string{"tags": {"tag1", "tag2"}},
		intSliceMap:    map[string][]int64{"values": {10, 20, 30}},
		floatSliceMap:  map[string][]float64{"rates": {0.1, 0.2}},
	}

	// Convert to ContextReader and parse
	config := Parse[Config](V3(cmdReader))

	// Verify the values were correctly parsed
	assert.Equal(t, "test-config", config.Name)
	assert.Equal(t, 25, config.Age)
	assert.Equal(t, true, config.Enabled)
	assert.Equal(t, 49.99, config.Price)
	assert.Equal(t, 10*time.Minute, config.Timeout)
	assert.Equal(t, []string{"tag1", "tag2"}, config.Tags)
	assert.Equal(t, []int{10, 20, 30}, config.Values)
	assert.Equal(t, []float64{0.1, 0.2}, config.Rate)
}
