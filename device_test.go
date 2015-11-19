// +build linux

package udev

import (
	"fmt"
	"runtime"
	"testing"
)

func ExampleDevice() {

	// Create Udev
	u := Udev{}

	// Create new Device based on subsystem and sysname
	d := u.NewDeviceFromSubsystemSysname("mem", "zero")

	// Extract information
	fmt.Printf("Syspath:%v\n", d.Syspath())
	fmt.Printf("Devpath:%v\n", d.Devpath())
	fmt.Printf("Devnode:%v\n", d.Devnode())
	fmt.Printf("Subsystem:%v\n", d.Subsystem())
	fmt.Printf("Devtype:%v\n", d.Devtype())
	fmt.Printf("Sysnum:%v\n", d.Sysnum())
	fmt.Printf("IsInitialized:%v\n", d.IsInitialized())
	fmt.Printf("Driver:%v\n", d.Driver())

	// Use one of the iterators
	it := d.PropertyIterator()
	it.Each(func(item interface{}) {
		kv := item.([]string)
		_ = fmt.Sprintf("Property:%v=%v\n", kv[0], kv[1])
	})
	// Output:
	// Syspath:/sys/devices/virtual/mem/zero
	// Devpath:/devices/virtual/mem/zero
	// Devnode:/dev/zero
	// Subsystem:mem
	// Devtype:
	// Sysnum:
	// IsInitialized:true
	// Driver:
}

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
	if d.PropertyValue("SUBSYSTEM") != "mem" {
		t.Fail()
	}
	if !d.IsInitialized() {
		t.Fail()
	}
	if d.SysattrValue("subsystem") != "mem" {
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
	it := d.PropertyIterator()
	it.Each(func(item interface{}) {
		_ = item.([]string)
	})
}

func TestDeviceGC(t *testing.T) {
	runtime.GC()
}
