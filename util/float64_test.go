package util

import "testing"

func TestFormatFloatStrDownByMultiplier(t *testing.T) {
	_, str1 := FormatFloatStrUpByMultiplier(4473.033, 0.02, 2)
	if str1 != "4473.04" {
		t.Errorf("FormatFloatStrDownByMultiplier failed, expect: 4473.03, got: %s", str1)
	}
	_, str2 := FormatFloatStrDownByMultiplier(4473.033, 0.02, 2)
	if str2 != "4473.02" {
		t.Errorf("FormatFloatStrDownByMultiplier failed, expect: 4473.03, got: %s", str2)
	}
}
