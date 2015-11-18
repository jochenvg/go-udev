// +build linux

package udev

import (
	"fmt"
	"testing"
)

func ExampleUdev_NewDeviceFromDevnum() {
	u := Udev{}
	d := u.NewDeviceFromDevnum('c', MkDev(1, 8))
	fmt.Println(d.Syspath())
	// Output:
	// /sys/devices/virtual/mem/random
}

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

func ExampleUdev_NewDeviceFromSyspath() {
	u := Udev{}
	d := u.NewDeviceFromSyspath("/sys/devices/virtual/mem/random")
	fmt.Println(d.Syspath())
	// Output:
	// /sys/devices/virtual/mem/random
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

func ExampleUdev_NewDeviceFromSubsystemSysname() {
	u := Udev{}
	d := u.NewDeviceFromSubsystemSysname("mem", "random")
	fmt.Println(d.Syspath())
	// Output:
	// /sys/devices/virtual/mem/random
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

func ExampleUdev_NewDeviceFromDeviceID() {
	u := Udev{}
	d := u.NewDeviceFromDeviceID("c1:8")
	fmt.Println(d.Syspath())
	// Output:
	// /sys/devices/virtual/mem/random
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

func ExampleUdev_NewMonitorFromNetlink() {
	u := Udev{}
	_ = u.NewMonitorFromNetlink("udev")
}

func TestNewMonitorFromNetlink(t *testing.T) {
	u := Udev{}
	_ = u.NewMonitorFromNetlink("udev")
}

func ExampleUdev_NewEnumerate() {
	u := Udev{}
	_ = u.NewEnumerate()
}

func TestNewEnumerate(t *testing.T) {
	u := Udev{}
	_ = u.NewEnumerate()
}
