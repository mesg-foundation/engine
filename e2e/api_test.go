package main

import (
	"testing"
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
}
