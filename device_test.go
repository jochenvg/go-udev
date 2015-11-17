// +build linux

// Package udev provides a cgo wrapper around the libudev C library
package udev

import (
	"runtime"
	"testing"
)

func TestDeviceZero(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromDeviceID("c1:5")
	if d.Subsystem() != "mem" {
		t.Fail()
	}
	if d.Syspath() != "/sys/devices/virtual/mem/zero" {
		t.Fail()
	}
	if d.Devnode() != "/dev/zero" {
		t.Fail()
	}
	p, e := d.PropertyValue("SUBSYSTEM")
	if e != nil || p != "mem" {
		t.Fail()
	}
	if !d.IsInitialized() {
		t.Fail()
	}
	s, e := d.SysattrValue("subsystem")
	if e != nil || s != "mem" {
		t.Fail()
	}
	// Device should have Properties
	properties := d.Properties()
	if len(properties) == 0 {
		t.Fail()
	}
	// Device should have Sysattrs
	sysattrs := d.Sysattrs()
	if len(sysattrs) == 0 {
		t.Fail()
	}
}

func TestDeviceRandom(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromDeviceID("c1:8")
	if d.Subsystem() != "mem" {
		t.Fail()
	}
	if d.Syspath() != "/sys/devices/virtual/mem/random" {
		t.Fail()
	}
	if d.Devnode() != "/dev/random" {
		t.Fail()
	}
	p, e := d.PropertyValue("SUBSYSTEM")
	if e != nil || p != "mem" {
		t.Fail()
	}
	if !d.IsInitialized() {
		t.Fail()
	}
	s, e := d.SysattrValue("subsystem")
	if e != nil || s != "mem" {
		t.Fail()
	}
	// Device should have Properties
	properties := d.Properties()
	if len(properties) == 0 {
		t.Fail()
	}
	// Device should have Sysattrs
	sysattrs := d.Sysattrs()
	if len(sysattrs) == 0 {
		t.Fail()
	}
}

func TestDeviceGC(t *testing.T) {
	runtime.GC()
}
