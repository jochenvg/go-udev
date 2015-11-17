// +build linux

// Package udev provides a cgo wrapper around the libudev C library
package udev

import (
	"runtime"
	"testing"
)

func TestEnumerateDeviceSyspaths(t *testing.T) {
	u := Udev{}
	e := u.NewEnumerate()
	dsp, err := e.DeviceSyspaths()
	if err != nil {
		t.Fail()
	}
	if len(dsp) <= 0 {
		t.Fail()
	}
}

func TestEnumerateSubsystemSyspaths(t *testing.T) {
	u := Udev{}
	e := u.NewEnumerate()
	ssp, err := e.SubsystemSyspaths()
	if err != nil {
		t.Fail()
	}
	if len(ssp) == 0 {
		t.Fail()
	}
}

func TestEnumerateDevicesWithFilter(t *testing.T) {
	u := Udev{}
	e := u.NewEnumerate()
	e.AddMatchSubsystem("block")
	e.AddMatchIsInitialized()
	e.AddNomatchSubsystem("mem")
	e.AddMatchProperty("ID_TYPE", "disk")
	e.AddMatchSysattr("partition", "1")
	e.AddMatchTag("systemd")
	//	e.AddMatchProperty("DEVTYPE", "partition")
	ds, err := e.Devices()
	if err != nil {
		t.Fail()
	}
	if len(ds) == 0 {
		t.Fail()
	}
	for k, d := range ds {
		if k != d.Syspath() {
			t.Fail()
		}
		if d.Subsystem() != "block" {
			t.Fail()
		}
		if !d.IsInitialized() {
			t.Fail()
		}
		if d.PropertyValue("ID_TYPE") != "disk" {
			t.Fail()
		}
		if d.SysattrValue("partition") != "1" {
			t.Fail()
		}
		if !d.HasTag("systemd") {
			t.Fail()
		}

		parent := d.Parent()
		parent2 := d.ParentWithSubsystemDevtype("block", "disk")
		if parent.Syspath() != parent2.Syspath() {
			t.Fail()
		}

	}
}

func TestEnumerateGC(t *testing.T) {
	runtime.GC()
}
