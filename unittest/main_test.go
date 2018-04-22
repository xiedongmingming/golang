package main

import "testing"

func TestTestFunc(t *testing.T) {

	r := TestFunc()

	if r != 3 {
		t.Errorf("TestFunc() failed. Got %d, expected 3.", r)
	}
	
}