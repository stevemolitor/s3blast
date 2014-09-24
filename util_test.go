package main

import (
        "testing"
        "sort"
)

func equal(a []string, b []string) bool {
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

func TestUniq(t *testing.T) {
        a := []string{"a", "b", "a", "c", "b"}
        actual := uniq(a)
        sort.Strings(actual)

        expected := []string{"a", "b", "c"}
        if (!equal(expected, actual)) {
                t.Errorf("uniq failed, actual results:", actual, "len:", len(actual))
        }
}
