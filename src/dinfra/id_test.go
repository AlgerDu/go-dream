package dinfra

import "testing"

func TestID_Zero(t *testing.T) {
	var id ID
	t.Log(id.Base32())
	t.Log(id.Base58())
}
