package utils

import "testing"

func TestToPtr(t *testing.T) {
	// Test for int
	i := 5
	ptri := ToPtr(i)
	if *ptri != i {
		t.Errorf("Expected %d, got %d", i, *ptri)
	}

	i = 0
	ptri = ToPtr(i)
	if ptri != nil {
		t.Errorf("Expected nil, got %d", *ptri)
	}

	// Test for string
	s := "hello"
	ptrs := ToPtr(s)
	if *ptrs != s {
		t.Errorf("Expected %s, got %s", s, *ptrs)
	}

	s = ""
	ptrs = ToPtr(s)
	if ptrs != nil {
		t.Errorf("Expected nil, got %s", *ptrs)
	}
}
