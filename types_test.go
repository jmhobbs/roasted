package main

import (
	"testing"
	"time"
)

func TestTemperatureValid(t *testing.T) {
	var temp Temperature

	temp = 0xFF00
	if temp.Valid() {
		t.Error("Expected Temperature to be invalid, but was valid.")
	}

	temp = 0x00AF
	if !temp.Valid() {
		t.Error("Expected Temperature to be valid, but was invalid.")
	}
}

func TestTimerValue(t *testing.T) {
	tests := []struct {
		Timer    Timer
		Expected time.Duration
	}{
		{
			0x0F,
			time.Minute + time.Second*30,
		},
		{
			0x3B,
			time.Minute*5 + time.Second*54,
		},
	}

	for _, test := range tests {
		actual := test.Timer.Value()
		if actual != test.Expected {
			t.Errorf("Value does not match.\nexpected: %v\n  actual: %v", test.Expected, actual)
		}
	}
}
