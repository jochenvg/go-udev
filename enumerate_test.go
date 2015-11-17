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
	if err != nil || len(ds) == 0 {
		t.Fail()
	}
	for i := range ds {
		if ds[i].Subsystem() != "block" {
			t.Fail()
		}
		if !ds[i].IsInitialized() {
			t.Fail()
		}
		value, e := ds[i].PropertyValue("ID_TYPE")
		if e != nil || value != "disk" {
			t.Fail()
		}
		value, e = ds[i].SysattrValue("partition")
		if e != nil || value != "1" {
			t.Fail()
		}
		if !ds[i].HasTag("systemd") {
			t.Fail()
		}

		parent := ds[i].Parent()
		if e != nil {
			t.Fail()
		}

		parent2 := ds[i].ParentWithSubsystemDevtype("block", "disk")
		if parent.Syspath() != parent2.Syspath() {
			t.Fail()
		}

	}
}

func TestEnumerateGC(t *testing.T) {
	runtime.GC()
}
