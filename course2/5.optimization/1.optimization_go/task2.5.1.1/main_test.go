package main

import (
	"testing"
)

func TestMapSlice(t *testing.T) {
	m := NewHashMapSlice(20)
	m.Set("Hello", "World")

	result, ok := m.Get("Hello")
	expected := "World"

	if ok != true {
		t.Errorf("got false, expected true")
	}

	if result != expected {
		t.Errorf("Got %v, expected %v", result, expected)
	}

	_, ok = m.Get("true")

	if ok != false {
		t.Errorf("Expected false, got %v", ok)
	}
}

func TestMapList(t *testing.T) {
	m := NewHashMapList(20)
	m.Set("Hello", "World")

	result, ok := m.Get("Hello")
	expected := "World"

	if ok != true {
		t.Errorf("got false, expected true")
	}

	if result != expected {
		t.Errorf("Got %v, expected %v", result, expected)
	}

	_, ok = m.Get("true")

	if ok != false {
		t.Errorf("Expected false, got %v", ok)
	}
}

func TestMapsWithZeroCap(t *testing.T) {
	m1 := NewHashMapList(0)
	m1.Set("Hello", "World")
	m2 := NewHashMapSlice(0)
	m2.Set("Hello", "World")

	result1, ok1 := m1.Get("Hello")
	result2, ok2 := m2.Get("Hello")

	if result1 != nil && result2 != nil {
		t.Errorf("Wrong behavior with 0 cap maps")
	}

	if ok1 != false && ok2 != false {
		t.Errorf("Wrong behavior with 0 cap maps")
	}
}
