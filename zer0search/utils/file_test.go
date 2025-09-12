package utils

import "testing"

func TestFileExists(t *testing.T) {
	testCase := []struct {
		filename string
		flag bool
	} {
		{"file.go", true},
		{"a.txt", false},
	}

	for _, test := range testCase {
		got := FileExists(test.filename)
		if test.flag != got {
			t.Errorf("FileExists(%s) = %v want %v", test.filename, got, test.flag)
		}
	}
}
