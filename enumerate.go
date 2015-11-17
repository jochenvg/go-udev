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

// Private structure wrapping a udev enumerate object
type enumerate struct {
	ptr *C.struct_udev_enumerate
	u   *Udev
}

// Lock the udev context
func (e *enumerate) lock() {
	e.u.m.Lock()
}

// Unlock the udev context
func (e *enumerate) unlock() {
	e.u.m.Unlock()
}

func enumerateUnref(e *enumerate) {
	C.udev_enumerate_unref(e.ptr)
}

func (e *enumerate) AddMatchSubsystem(subsystem string) (err error) {
	e.lock()
	defer e.unlock()
	s := C.CString(subsystem)
	defer freeCharPtr(s)
	if C.udev_enumerate_add_match_subsystem(e.ptr, s) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_subsystem failed")
	}
	return
}

func (e *enumerate) AddNomatchSubsystem(subsystem string) (err error) {
	e.lock()
	defer e.unlock()
	s := C.CString(subsystem)
	defer freeCharPtr(s)
	if C.udev_enumerate_add_nomatch_subsystem(e.ptr, s) != 0 {
		err = errors.New("udev: udev_enumerate_add_nomatch_subsystem failed")
	}
	return
}

func (e *enumerate) AddMatchSysattr(sysattr, value string) (err error) {
	e.lock()
	defer e.unlock()
	s, v := C.CString(sysattr), C.CString(value)
	defer freeCharPtr(s)
	defer freeCharPtr(v)
	if C.udev_enumerate_add_match_sysattr(e.ptr, s, v) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_sysattr failed")
	}
	return
}

func (e *enumerate) AddNomatchSysattr(sysattr, value string) (err error) {
	e.lock()
	defer e.unlock()
	s, v := C.CString(sysattr), C.CString(value)
	defer freeCharPtr(s)
	defer freeCharPtr(v)
	if C.udev_enumerate_add_nomatch_sysattr(e.ptr, s, v) != 0 {
		err = errors.New("udev: udev_enumerate_add_nomatch_sysattr failed")
	}
	return
}

func (e *enumerate) AddMatchProperty(property, value string) (err error) {
	e.lock()
	defer e.unlock()
	p, v := C.CString(property), C.CString(value)
	defer freeCharPtr(p)
	defer freeCharPtr(v)
	if C.udev_enumerate_add_match_property(e.ptr, p, v) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_property failed")
	}
	return
}

func (e *enumerate) AddMatchSysname(sysname string) (err error) {
	e.lock()
	defer e.unlock()
	s := C.CString(sysname)
	defer freeCharPtr(s)
	if C.udev_enumerate_add_match_sysname(e.ptr, s) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_sysname failed")
	}
	return
}

func (e *enumerate) AddMatchTag(tag string) (err error) {
	e.lock()
	defer e.unlock()
	t := C.CString(tag)
	defer freeCharPtr(t)
	if C.udev_enumerate_add_match_tag(e.ptr, t) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_tag failed")
	}
	return
}

func (e *enumerate) AddMatchParent(parent *device) (err error) {
	e.lock()
	defer e.unlock()
	if C.udev_enumerate_add_match_parent(e.ptr, parent.ptr) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_parent failed")
	}
	return
}

func (e *enumerate) AddMatchIsInitialized() (err error) {
	e.lock()
	defer e.unlock()
	if C.udev_enumerate_add_match_is_initialized(e.ptr) != 0 {
		err = errors.New("udev: udev_enumerate_add_match_is_initialized failed")
	}
	return
}

func (e *enumerate) AddSyspath(syspath string) (err error) {
	e.lock()
	defer e.unlock()
	s := C.CString(syspath)
	defer freeCharPtr(s)
	if C.udev_enumerate_add_syspath(e.ptr, s) != 0 {
		err = errors.New("udev: udev_enumerate_add_syspath failed")
	}
	return
}

func (e *enumerate) DeviceSyspaths() (s Set, err error) {
	e.lock()
	defer e.unlock()
	if C.udev_enumerate_scan_devices(e.ptr) < 0 {
		err = errors.New("udev: udev_enumerate_scan_devices failed")
	} else {
		s = make(Set)
		s.addFromListEntry(C.udev_enumerate_get_list_entry(e.ptr))
	}
	return
}

func (e *enumerate) SubsystemSyspaths() (s Set, err error) {
	e.lock()
	defer e.unlock()
	if C.udev_enumerate_scan_subsystems(e.ptr) < 0 {
		err = errors.New("udev: udev_enumerate_scan_subsystems failed")
	} else {
		s = make(Set)
		s.addFromListEntry(C.udev_enumerate_get_list_entry(e.ptr))
	}
	return
}

func (e *enumerate) Devices() (m DeviceMap, err error) {
	e.lock()
	defer e.unlock()
	if C.udev_enumerate_scan_devices(e.ptr) < 0 {
		err = errors.New("udev: udev_enumerate_scan_devices failed")
	} else {
		m = make(DeviceMap)
		for l := C.udev_enumerate_get_list_entry(e.ptr); l != nil; l = C.udev_list_entry_get_next(l) {
			s := C.udev_list_entry_get_name(l)
			m[C.GoString(s)] = e.u.newDevice(C.udev_device_new_from_syspath(e.u.ptr, s))
		}
	}
	return
}
