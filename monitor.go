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
import (
	"errors"

	"golang.org/x/sys/unix"
)

// Private structure wrapping a udev monitor object
type monitor struct {
	ptr *C.struct_udev_monitor
	u   *Udev
}

const (
	maxEpollEvents = 32
	epollTimeout   = 1000
)

// Lock the udev context
func (m *monitor) lock() {
	m.u.m.Lock()
}

// Unlock the udev context
func (m *monitor) unlock() {
	m.u.m.Unlock()
}

func monitorUnref(m *monitor) {
	C.udev_monitor_unref(m.ptr)
}

func (m *monitor) SetReceiveBufferSize(size int) (err error) {
	m.lock()
	defer m.unlock()
	if C.udev_monitor_set_receive_buffer_size(m.ptr, (C.int)(size)) != 0 {
		err = errors.New("udev: udev_monitor_set_receive_buffer_size failed")
	}
	return
}

func (m *monitor) FilterAddMatchSubsystemDevtype(subsystem, devtype string) (err error) {
	m.lock()
	defer m.unlock()
	s, d := C.CString(subsystem), C.CString(devtype)
	defer freeCharPtr(s)
	defer freeCharPtr(d)
	if C.udev_monitor_filter_add_match_subsystem_devtype(m.ptr, s, d) != 0 {
		err = errors.New("udev: udev_monitor_filter_add_match_subsystem_devtype failed")
	}
	return
}

func (m *monitor) FilterAddMatchTag(tag string) (err error) {
	m.lock()
	defer m.unlock()
	t := C.CString(tag)
	defer freeCharPtr(t)
	if C.udev_monitor_filter_add_match_tag(m.ptr, t) != 0 {
		err = errors.New("udev: udev_monitor_filter_add_match_tag failed")
	}
	return
}

func (m *monitor) FilterUpdate() (err error) {
	m.lock()
	defer m.unlock()
	if C.udev_monitor_filter_update(m.ptr) != 0 {
		err = errors.New("udev: udev_monitor_filter_update failed")
	}
	return
}

func (m *monitor) FilterRemove() (err error) {
	m.lock()
	defer m.unlock()
	if C.udev_monitor_filter_remove(m.ptr) != 0 {
		err = errors.New("udev: udev_monitor_filter_remove failed")
	}
	return
}

func (m *monitor) receiveDevice() (d *device) {
	m.lock()
	defer m.unlock()
	return m.u.newDevice(C.udev_monitor_receive_device(m.ptr))
}

func (m *monitor) DeviceChan(done <-chan struct{}) (<-chan *device, error) {

	var event unix.EpollEvent
	var events [maxEpollEvents]unix.EpollEvent

	// Lock the context
	m.lock()
	defer m.unlock()

	// Enable receiving
	if C.udev_monitor_enable_receiving(m.ptr) != 0 {
		return nil, errors.New("udev: udev_monitor_enable_receiving failed")
	}

	// Set the fd to non-blocking
	fd := C.udev_monitor_get_fd(m.ptr)
	if e := unix.SetNonblock(int(fd), true); e != nil {
		return nil, errors.New("udev: unix.SetNonblock failed")
	}

	// Create an epoll fd
	epfd, e := unix.EpollCreate1(0)
	if e != nil {
		return nil, errors.New("udev: unix.EpollCreate1 failed")
	}

	// Add the fd to the epoll fd
	event.Events = unix.EPOLLIN | unix.EPOLLET
	event.Fd = int32(fd)
	if e = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(fd), &event); e != nil {
		return nil, errors.New("udev: unix.EpollCtl failed")
	}

	// Create the channel
	ch := make(chan *device)

	// Create goroutine to epoll the fd
	go func(done <-chan struct{}, fd int32) {
		// Close the epoll fd when goroutine exits
		defer unix.Close(epfd)
		// Close the channel when goroutine exits
		defer close(ch)
		// Loop forever
		for {
			// Poll the file descriptor
			nevents, e := unix.EpollWait(epfd, events[:], epollTimeout)
			if e != nil {
				return
			}
			// Process events
			for ev := 0; ev < nevents; ev++ {
				if events[ev].Fd == fd {
					if (events[ev].Events & unix.EPOLLIN) != 0 {
						for d := m.receiveDevice(); d != nil; d = m.receiveDevice() {
							ch <- d
						}
					}
				}
			}
			// Check for done signal
			select {
			case <-done:
				return
			default:
			}
		}
	}(done, int32(fd))

	return ch, nil
}
