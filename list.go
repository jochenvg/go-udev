// +build linux,cgo

package udev

/*
  #cgo LDFLAGS: -ludev
  #include <libudev.h>
  #include <linux/types.h>
  #include <stdlib.h>
	#include <linux/kdev_t.h>
*/
import "C"

// Set represents a set of strings
type Set map[string]struct{}

func (s Set) add(key string) {
	s[key] = struct{}{}
}

func (s Set) addFromListEntry(l *C.struct_udev_list_entry) {
	for ; l != nil; l = C.udev_list_entry_get_next(l) {
		s.add(C.GoString(C.udev_list_entry_get_name(l)))
	}
}

// Map represents a key/value map
type Map map[string]string

func (m Map) addFromListEntry(l *C.struct_udev_list_entry) {
	for ; l != nil; l = C.udev_list_entry_get_next(l) {
		m[C.GoString(C.udev_list_entry_get_name(l))] = C.GoString(C.udev_list_entry_get_value(l))
	}
}

// DeviceMap is a map from syspaths to device
type DeviceMap map[string]*Device

const (
	listCapacity       = 100
	deviceListCapacity = 100
)

// List is a list of strings
type List []string

func (list List) addFromListEntry(l *C.struct_udev_list_entry) List {
	for ; l != nil; l = C.udev_list_entry_get_next(l) {
		list = append(list, C.GoString(C.udev_list_entry_get_name(l)))
	}
	return list
}

// DeviceList is a list of Device pointers
type DeviceList [](*Device)
