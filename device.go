// +build linux,cgo

// Package udev provides a cgo wrapper around the libudev C library
package udev

/*
  #cgo LDFLAGS: -ludev
  #include <libudev.h>
  #include <linux/types.h>
  #include <stdlib.h>
	#include <linux/kdev_t.h>
*/
import "C"
import "errors"

// Private structure wrapping a udev device object
type device struct {
	ptr *C.struct_udev_device
	u   *Udev
}

// Lock the udev context
func (d *device) lock() {
	d.u.m.Lock()
}

// Unlock the udev context
func (d *device) unlock() {
	d.u.m.Unlock()
}

func deviceUnref(d *device) {
	C.udev_device_unref(d.ptr)
}

func (d *device) Parent() *device {
	d.lock()
	defer d.unlock()
	ptr := C.udev_device_get_parent(d.ptr)
	if ptr != nil {
		C.udev_device_ref(ptr)
	}
	return d.u.newDevice(ptr)
}

func (d *device) ParentWithSubsystemDevtype(subsystem, devtype string) *device {
	d.lock()
	defer d.unlock()
	ss, dt := C.CString(subsystem), C.CString(devtype)
	defer freeCharPtr(ss)
	defer freeCharPtr(dt)
	ptr := C.udev_device_get_parent_with_subsystem_devtype(d.ptr, ss, dt)
	if ptr != nil {
		C.udev_device_ref(ptr)
	}
	return d.u.newDevice(ptr)
}

func (d *device) DevPath() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_devpath(d.ptr))
}

func (d *device) Subsystem() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_subsystem(d.ptr))
}

func (d *device) DevType() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_devtype(d.ptr))
}

func (d *device) Syspath() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_syspath(d.ptr))
}

func (d *device) Sysnum() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_sysnum(d.ptr))
}

func (d *device) Devnode() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_devnode(d.ptr))
}

func (d *device) IsInitialized() bool {
	d.lock()
	defer d.unlock()
	return C.udev_device_get_is_initialized(d.ptr) != 0
}

func (d *device) DevLinks() (r Set) {
	d.lock()
	defer d.unlock()
	r = make(Set)
	r.addFromListEntry(C.udev_device_get_devlinks_list_entry(d.ptr))
	return
}

func (d *device) Properties() (r Map) {
	d.lock()
	defer d.unlock()
	r = make(Map)
	r.addFromListEntry(C.udev_device_get_properties_list_entry(d.ptr))
	return
}

func (d *device) Tags() (r Set) {
	d.lock()
	defer d.unlock()
	r = make(Set)
	r.addFromListEntry(C.udev_device_get_tags_list_entry(d.ptr))
	return
}

func (d *device) Sysattrs() (r Set) {
	d.lock()
	defer d.unlock()
	r = make(Set)
	r.addFromListEntry(C.udev_device_get_sysattr_list_entry(d.ptr))
	return
}

func (d *device) PropertyValue(key string) string {
	d.lock()
	defer d.unlock()
	k := C.CString(key)
	defer freeCharPtr(k)
	return C.GoString(C.udev_device_get_property_value(d.ptr, k))
}

func (d *device) Driver() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_driver(d.ptr))
}

func (d *device) Devnum() Devnum {
	d.lock()
	defer d.unlock()
	return Devnum{C.udev_device_get_devnum(d.ptr)}
}

func (d *device) Action() string {
	d.lock()
	defer d.unlock()
	return C.GoString(C.udev_device_get_action(d.ptr))
}

func (d *device) Seqnum() uint64 {
	d.lock()
	defer d.unlock()
	return uint64(C.udev_device_get_seqnum(d.ptr))
}

func (d *device) UsecSinceInitialized() uint64 {
	d.lock()
	defer d.unlock()
	return uint64(C.udev_device_get_usec_since_initialized(d.ptr))
}

func (d *device) SysattrValue(sysattr string) string {
	d.lock()
	defer d.unlock()
	s := C.CString(sysattr)
	defer freeCharPtr(s)
	return C.GoString(C.udev_device_get_sysattr_value(d.ptr, s))
}

func (d *device) SetSysattrValue(sysattr, value string) (err error) {
	d.lock()
	defer d.unlock()
	sa, val := C.CString(sysattr), C.CString(value)
	defer freeCharPtr(sa)
	defer freeCharPtr(val)
	if C.udev_device_set_sysattr_value(d.ptr, sa, val) < 0 {
		err = errors.New("udev: udev_device_set_sysattr_value failed")
	}
	return
}

func (d *device) HasTag(tag string) bool {
	d.lock()
	defer d.unlock()
	t := C.CString(tag)
	defer freeCharPtr(t)
	return C.udev_device_has_tag(d.ptr, t) != 0
}
