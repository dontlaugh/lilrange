package lilrange

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type testCase struct {
		input    string
		expected *Range
	}
	var invalidInputs []testCase = []testCase{
		// invalid inputs, so we expect nil
		{"0300", nil},
		{"0300-000", nil},
		{"300-000", nil},
		{"-000", nil},
		{"", nil},
		{"-", nil},
		{"-0300", nil},
		{"-0", nil},
		{"0030-0", nil},
	}

	for _, c := range invalidInputs {
		actual, _ := Parse(c.input)
		assert.EqualValues(t, c.expected, actual)
	}

	type durationTest struct {
		input    string
		expected time.Duration
	}
	var durationTests []durationTest = []durationTest{
		// valid inputs, and we assert on the durations we compute
		{"0100-0200", 1 * time.Hour},
		{"0000-0002", 2 * time.Minute},
		{"2359-0000", 1 * time.Minute},
		{"2359-0001", 2 * time.Minute},
		{"2259-0001", 62 * time.Minute},
	}

	for _, dt := range durationTests {
		actual, err := Parse(dt.input)
		if err != nil {
			t.Errorf("duration test fail: %v", err)
		}
		assert.EqualValues(t, dt.expected, actual.Duration)
	}
}

func TestExtractAndValidate(t *testing.T) {
	type testCase struct {
		input    string
		expected []int
	}
	var cases = []testCase{
		// invalid cases
		{"0", []int{-1, -1}},
		{"223", []int{-1, -1}},
		{"2460", []int{-1, -1}},
		{"2460", []int{-1, -1}},
		{"0260", []int{-1, -1}},
		{"9901", []int{-1, -1}},
		{"12.1", []int{-1, -1}},
		{"12:1", []int{-1, -1}},
		{"1:21", []int{-1, -1}},
		{"1:21", []int{-1, -1}},
		{"1/21", []int{-1, -1}},
		// valid cases
		{"2359", []int{23, 59}},
		{"0000", []int{0, 0}},
		{"0001", []int{0, 1}},
		{"1111", []int{11, 11}},
		{"0223", []int{2, 23}},
	}

	for _, c := range cases {
		hr, min, _ := extractAndValidate(c.input)
		assert.EqualValues(t, c.expected, []int{hr, min})
	}
}

func TestCalculateDurationMinutes(t *testing.T) {

	type testCase struct {
		// inputs: startHr, startMin, endHr, endMin
		inputs   []int
		expected int
	}
	var cases = []testCase{
		{[]int{0, 0, 1, 0}, 60},
		{[]int{1, 0, 2, 0}, 60},
		{[]int{23, 0, 1, 0}, 120},
		{[]int{23, 15, 1, 0}, 105},
		{[]int{23, 15, 1, 10}, 115},
	}

	for _, c := range cases {
		dur, _ := CalculateDurationMinutes(c.inputs[0], c.inputs[1],
			c.inputs[2], c.inputs[3])
		assert.EqualValues(t, c.expected, dur)
	}
	// Assert that invalid inputs panic.
	assert.Panics(t, func() { CalculateDurationMinutes(-1, 0, 0, 0) })
	assert.Panics(t, func() { CalculateDurationMinutes(25, 0, 0, 0) })
	assert.Panics(t, func() { CalculateDurationMinutes(0, 0, 25, 0) })
	assert.Panics(t, func() { CalculateDurationMinutes(0, 69, 0, 0) })
	assert.Panics(t, func() { CalculateDurationMinutes(0, 0, 0, 69) })
}
