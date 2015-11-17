// +build linux

// Package udev provides a cgo wrapper around the libudev C library
package udev

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestMonitorDeviceChan(t *testing.T) {
	u := Udev{}
	m := u.NewMonitorFromNetlink("udev")
	m.FilterAddMatchSubsystemDevtype("block", "disk")
	m.FilterAddMatchTag("systemd")
	done := make(chan struct{})
	ch, e := m.DeviceChan(done)
	if e != nil {
		t.Fail()
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("Started listening on channel")
		for d := range ch {
			fmt.Println(d.Syspath(), d.Action())
		}
		fmt.Println("Channel closed")
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to update filter")
		<-time.After(2 * time.Second)
		fmt.Println("Removing filter")
		m.FilterRemove()
		fmt.Println("Updating filter")
		m.FilterUpdate()
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to signal done")
		<-time.After(4 * time.Second)
		fmt.Println("Signalling done")
		close(done)
		wg.Done()
	}()
	wg.Wait()

}

func TestMonitorGC(t *testing.T) {
	runtime.GC()
}
