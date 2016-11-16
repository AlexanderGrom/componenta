package config

import (
	"os"
	"reflect"
	"testing"
)

var cfg Config

func TestCommon(t *testing.T) {
	os.Setenv("TESTVARNAME", "componenta")
	var err error
	cfg, err = Use("./test.conf")

	if err != nil {
		t.Error(err)
	}
}

func TestKey1(t *testing.T) {
	expect := `value`
	result := cfg.GetString("key1")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey2(t *testing.T) {
	expect := `value2`
	result := cfg.GetString("key2")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey3(t *testing.T) {
	var expect int64 = 11
	result := cfg.GetInt("key3")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey4(t *testing.T) {
	var expect float64 = 3.14
	result := cfg.GetFloat("key4")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey4.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey5(t *testing.T) {
	expect := true
	result := cfg.GetBool("key5")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey5.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey6(t *testing.T) {
	expect := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := cfg.GetInts("key6")

	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in TestKey6.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey7(t *testing.T) {
	expect := []bool{true, false}
	result := cfg.GetBools("key7")

	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in TestKey7.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey8(t *testing.T) {
	expect := []string{"foo", "bar"}
	result := cfg.GetStrings("key8")

	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in TestKey8.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey9(t *testing.T) {
	expect := ``
	result := cfg.GetString("key9")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey9.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey10(t *testing.T) {
	var expect int64 = 0
	result := cfg.GetInt("key10")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey10.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey11(t *testing.T) {
	expect := false
	result := cfg.GetBool("key11")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey11.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey12(t *testing.T) {
	expect := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	result := cfg.GetInts("key12")

	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in TestKey12.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey13(t *testing.T) {
	expect := `value
         value
         value`
	result := cfg.GetString("key13")
	if expect != result {
		t.Errorf("Expect result to equal in TestKey13.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey14(t *testing.T) {
	expect := `value["key"]`
	result := cfg.GetString("key14")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey14.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey15(t *testing.T) {
	expect := `componenta`
	result := cfg.GetString("key15")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey15.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey16(t *testing.T) {
	expect := `value`
	result := cfg.GetString("key16")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey16.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey17(t *testing.T) {
	expect := []int64{1, 2, 3, 4, 5, 6}
	result := cfg.GetInts("key17")

	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in TestKey17.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey18(t *testing.T) {
	expect := `value`
	result := cfg.GetString("key18")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey18.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey19(t *testing.T) {
	expect := `value`
	result := cfg.GetString("key19")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey19.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey20(t *testing.T) {
	expect := `componenta/src`
	result := cfg.GetString("key20")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey20.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey21(t *testing.T) {
	expect := `${TESTVARNAME/src`
	result := cfg.GetString("key21")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey21.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey22(t *testing.T) {
	expect := `/var/componenta/src`
	result := cfg.GetString("key22")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey22.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestKey23(t *testing.T) {
	expect := `${TESTVARNAME}`
	result := cfg.GetString("key23")

	if expect != result {
		t.Errorf("Expect result to equal in TestKey23.\nResult: %s\nExpect: %s", result, expect)
	}
}

func DataEqual(aa, bb interface{}) bool {
	a := InterfaceSlice(aa)
	b := InterfaceSlice(bb)
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func InterfaceSlice(x interface{}) []interface{} {
	s := reflect.ValueOf(x)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	a := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		a[i] = s.Index(i).Interface()
	}
	return a
}
