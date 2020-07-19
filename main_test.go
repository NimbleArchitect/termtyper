package main

import "testing"

func TestValidCmdType(t *testing.T) {

	out, sep := validCmdType(" bASh ")
	if out != "bash" {
		t.Errorf("validCmdType string return incorrect value, got: %s, want: %s.", out, "bash")
	}
	if sep != " \\" {
		t.Errorf("validCmdType seperator return incorrect value, got: %s, want: %s.", out, " \\")
	}
}
