package tools

import "testing"

func TestHash(t *testing.T) {
	t.Log(Hash("123123123"))
	t.Log(Hash("444"))
	t.Log(Hash("555"))

}
