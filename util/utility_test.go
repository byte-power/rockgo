package util

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAnyTo(t *testing.T) {
	m := AnyMap{"1": 1}
	assert.Equal(t, m, AnyToAnyMap(m))
	assert.Equal(t, m, AnyToAnyMap(map[interface{}]interface{}{1: 1}))
	assert.Equal(t, "1", AnyToString(1))
	assert.Equal(t, "-1", AnyToString(-1))
	assert.Equal(t, "1.5", AnyToString(1.5))
	assert.Equal(t, "true", AnyToString(true))
	assert.Equal(t, fmt.Sprint(m), AnyToString(m))

	assert.Equal(t, int64(1), AnyToInt64(1.5))
	assert.Equal(t, int64(1), AnyToInt64(true))
	assert.Equal(t, int64(10), AnyToInt64("10"))
	assert.Equal(t, int64(10), AnyToInt64("10.5"))
	assert.Equal(t, int64(0), AnyToInt64("10a"))

	assert.Equal(t, float64(1), AnyToFloat64(1))
	assert.Equal(t, float64(1), AnyToFloat64("1"))
	assert.Equal(t, float64(0), AnyToFloat64("1a"))
	assert.Equal(t, float64(1), AnyToFloat64(true))

	assert.Equal(t, true, AnyToBool(1))
	assert.Equal(t, true, AnyToBool(1.5))
	assert.Equal(t, true, AnyToBool("1"))
	assert.Equal(t, true, AnyToBool("T"))
}

func TestFindInMap(t *testing.T) {
	m3 := AnyMap{"aaa": 1.5}
	arr := []int{1, 2}
	m := AnyMap{
		"a": AnyMap{"aa": m3},
		"b": arr,
	}
	assert.Equal(t, 1.5, FindInAnyMap(m, "a", "aa", "aaa"))
	assert.Equal(t, nil, FindInAnyMap(m, "a", "aa", "aab"))
	assert.Equal(t, arr, FindInAnyMap(m, "b").([]int))
}

func TestAnyArrayToMap(t *testing.T) {
	elements := []interface{}{"abc"}
	m := AnyArrayToMap(elements)
	if m != nil {
		t.Error("Should be nil if length of input < 2.")
	}
	elements2 := []interface{}{"abc", "bd", "cc"}
	m2 := AnyArrayToMap(elements2)
	if len(m2) != 1 {
		t.Error("Generated map length should 1.")
	}
}
