package util

import (
	"fmt"
	"sync"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAnyTo(t *testing.T) {
	m := StrMap{"1": 1}
	assert.Equal(t, m, AnyToStrMap(m))
	assert.Equal(t, m, AnyToStrMap(AnyMap{1: 1}))
	assert.Equal(t, AnyMap{"1": 1}, AnyToAnyMap(m))
	assert.Equal(t, "1", AnyToString(1))
	assert.Equal(t, "-1", AnyToString(-1))
	assert.Equal(t, "1.5", AnyToString(1.5))
	assert.Equal(t, "true", AnyToString(true))
	assert.Equal(t, fmt.Sprint(m), AnyToString(m))

	assert.Equal(t, int64(1), AnyToInt64(int8(1)))
	assert.Equal(t, int64(1), AnyToInt64(int16(1)))
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
	var sp *string
	assert.Equal(t, false, AnyToBool(sp))
	assert.Equal(t, int64(0), AnyToInt64(sp))
	assert.Equal(t, 0.0, AnyToFloat64(sp))
	assert.Empty(t, AnyToString(sp))
	assert.Nil(t, AnyToAnyMap(sp))
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

	ms := AnyToStrMap(m)
	assert.Equal(t, 1.5, FindInStrMap(ms, "a", "aa", "aaa"))
}

func TestFindInSyncMap(t *testing.T) {
	var m sync.Map
	m.Store("a", map[string]interface{}{"b": "c"})
	assert.Equal(t, "c", FindInSyncMap(&m, "a", "b"))
}

func TestAnyArrayToMap(t *testing.T) {
	elements := []interface{}{"abc"}
	m := AnyArrayToStrMap(elements)
	if m != nil {
		t.Error("Should be nil if length of input < 2.")
	}
	elements2 := []interface{}{"abc", "bd", "cc"}
	m2 := AnyArrayToStrMap(elements2)
	if len(m2) != 1 {
		t.Error("Generated map length should 1.")
	}
}
