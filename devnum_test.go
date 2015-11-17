// +build linux

// Package udev provides a cgo wrapper around the libudev C library
package udev

import "testing"

func TestDevnumMajorMinor(t *testing.T) {
	d := MkDev(1, 8)
	if d.Major() != 1 {
		t.Fail()
	}
	if d.Minor() != 8 {
		t.Fail()
	}
}
