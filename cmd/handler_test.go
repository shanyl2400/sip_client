package cmd

import "testing"

func TestParseArgs(t *testing.T) {
	pt := NewPrompt()
	out := pt.parseArgs([]string{
		"-n", "abc", "-p", "123456", "-x", "-w", "-t", "ttt",
	})
	t.Log(out)
}
