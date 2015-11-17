// +build linux

package udev

import "testing"

func TestNewDeviceFromDevnum(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromDevnum('c', MkDev(1, 8))
	if d.Devnum().Major() != 1 {
		t.Fail()
	}
	if d.Devnum().Minor() != 8 {
		t.Fail()
	}
	if d.Devpath() != "/devices/virtual/mem/random" {
		t.Fail()
	}
}

func TestNewDeviceFromDevnumNoClose(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromDevnum('c', MkDev(1, 8))
	if d.Devnum().Major() != 1 {
		t.Fail()
	}
	if d.Devnum().Minor() != 8 {
		t.Fail()
	}
	if d.Devpath() != "/devices/virtual/mem/random" {
		t.Fail()
	}
}

func TestNewDeviceFromSyspath(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromSyspath("/sys/devices/virtual/mem/random")
	if d.Devnum().Major() != 1 {
		t.Fail()
	}
	if d.Devnum().Minor() != 8 {
		t.Fail()
	}
	if d.Devpath() != "/devices/virtual/mem/random" {
		t.Fail()
	}
}

func TestNewDeviceFromSubsystemSysname(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromSubsystemSysname("mem", "random")
	if d.Devnum().Major() != 1 {
		t.Fail()
	}
	if d.Devnum().Minor() != 8 {
		t.Fail()
	}
	if d.Devpath() != "/devices/virtual/mem/random" {
		t.Fail()
	}
}

func TestNewDeviceFromDeviceID(t *testing.T) {
	u := Udev{}
	d := u.NewDeviceFromDeviceID("c1:8")
	if d.Devnum().Major() != 1 {
		t.Fail()
	}
	if d.Devnum().Minor() != 8 {
		t.Fail()
	}
	if d.Devpath() != "/devices/virtual/mem/random" {
		t.Fail()
	}
}

func TestNewMonitorFromNetlink(t *testing.T) {
	u := Udev{}
	_ = u.NewMonitorFromNetlink("udev")
}

func TestNewEnumerate(t *testing.T) {
	u := Udev{}
	_ = u.NewEnumerate()
}
