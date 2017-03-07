package crypt

import (
	"testing"
)

func TestCrypt(t *testing.T) {
	expect := "String"
	c, err := Encrypt(expect, "Secret_Key")

	if err != nil {
		t.Errorf("Encrypt Error:", err)
	}

	s, err := Decrypt(c, "Secret_Key")

	if err != nil {
		t.Errorf("Decrypt Error:", err)
	}

	if expect != s {
		t.Errorf("Expect result to equal in TestCrypt.\nResult: %s\nExpect: %s", s, expect)
	}
}
