package env

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Load("./test.conf")
	code := m.Run()
	os.Setenv("TEST_APP_HOME", "")
	os.Setenv("TEST_APP_HOST", "")
	os.Setenv("TEST_APP_PORT", "")
	os.Setenv("TEST_APP_URL", "")
	os.Exit(code)
}

func TestKey1(t *testing.T) {
	expect := `/var/app`
	result := os.Getenv("TEST_APP_HOME")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey2(t *testing.T) {
	expect := `localhost`
	result := os.Getenv("TEST_APP_HOST")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey3(t *testing.T) {
	expect := `1448`
	result := os.Getenv("TEST_APP_PORT")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey4(t *testing.T) {
	expect := `http://localhost:1448`
	result := os.Getenv("TEST_APP_URL")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey4.\nResult: %s\nExpect: %s", result, expect)
	}
}
