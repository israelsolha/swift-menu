package main

import (
	"testing"
)

func Test1(t *testing.T) {
	err := Testing(1)
	if err == nil {
		t.Errorf("Error")
	}
}

func Test2(t *testing.T) {
	err := Testing(2)
	if err == nil {
		t.Errorf("Error")
	}
}

func Test3(t *testing.T) {
	err := Testing(3)
	if err == nil {
		t.Errorf("Error")
	}
}

func Test4(t *testing.T) {
	err := Testing(4)
	if err == nil {
		t.Errorf("Error")
	}
}
