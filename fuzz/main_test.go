package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}

	for _, tc := range testcases {
		rev, _ := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		rev, revErr := Reverse(orig)
		if revErr != nil {
			return
		}

		doubleRev, doubleRevErr := Reverse(rev)
		if doubleRevErr != nil {
			return
		}

		// t.Logf("orig=%s, rev=%s, doubleRev=%s\n", orig, rev, doubleRev)

		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}

		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
