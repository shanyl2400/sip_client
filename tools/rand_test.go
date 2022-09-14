package tools

import "testing"

func TestRandString(t *testing.T) {
	t.Log(RandString())
	t.Log(RandString())
	t.Log(RandString())
}
