package util

import (
	"fmt"
	"testing"
)

func TestAnyTo(t *testing.T) {
	m := AnyMap{"1": 1}
	if AnyToAnyMap(m) == nil {
		t.Error("AnyToAnyMap(AnyMap{}) failed")
	}
	if AnyToAnyMap(map[interface{}]interface{}{1: 1}) == nil {
		t.Error("AnyToAnyMap(map[interface{}]interface{}{}) failed")
	}

	if AnyToString(1) != "1" {
		t.Error("AnyToString(1) != '1'")
	}
	if AnyToString(-1) != "-1" {
		t.Error("AnyToString(-1) != '-1'")
	}
	if AnyToString(1.5) != "1.5" {
		t.Error("AnyToString(1.5) != '1.5'")
	}
	if AnyToString(true) != "true" {
		t.Error("AnyToString(true) != 'true'")
	}
	if AnyToString(m) != fmt.Sprint(m) {
		t.Error("AnyToString(m) == fmt.Sprint(m)")
	}

	if AnyToInt64(1.5) != 1 {
		t.Error("AnyToInt64(1.5) != 1")
	}
	if AnyToInt64(true) != 1 {
		t.Error("AnyToInt64(true) != 1")
	}
	if AnyToInt64("10") != 10 {
		t.Error("AnyToInt64('10') != 10")
	}
	if AnyToInt64("10.5") != 10 {
		t.Error("AnyToInt64('10.5') != 10")
	}
	if AnyToInt64("10a") != 0 {
		t.Error("AnyToInt64('10a') != 0")
	}

	if AnyToFloat64(1) != 1.0 {
		t.Error("AnyToFloat64(1) != 1.0")
	}
	if AnyToFloat64("1") != 1.0 {
		t.Error("AnyToFloat64('1') != 1.0")
	}
	if AnyToFloat64("1a") != 0 {
		t.Error("AnyToFloat64('1a') != 1.0")
	}
	if AnyToFloat64(true) != 1.0 {
		t.Error("AnyToFloat64(true) != 1.0")
	}

	if AnyToBool(1) != true {
		t.Error("AnyToBool(1) != true")
	}
	if AnyToBool(1.5) != true {
		t.Error("AnyToBool(1.5) != true")
	}
	if AnyToBool("1") != true {
		t.Error("AnyToBool('1') != true")
	}
	if AnyToBool("T") != true {
		t.Error("AnyToBool('T') != true")
	}
}

func TestFindInMap(t *testing.T) {
	m3 := AnyMap{"aaa": 1.5}
	m := AnyMap{
		"a": AnyMap{
			"aa": m3,
		},
		"b": []int{1, 2},
	}
	if FindInAnyMap(m, "a", "aa", "aaa") != 1.5 {
		t.Error("FindInAnyMap 'a.aa.aaa' != 1.5")
	}
	if FindInAnyMap(m, "a", "aa", "aab") != nil {
		t.Error("FindInAnyMap 'a.aa.aab' != nil")
	}
	if _, ok := FindInAnyMap(m, "b").([]int); !ok {
		t.Error("FindInAnyMap 'b' != []int")
	}
}
