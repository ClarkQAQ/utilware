// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !aix,!darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package ipv4

import (
	"net"

	"utilware/dep/x/net/bpf"
	"utilware/dep/x/net/internal/socket"
)

func (so *sockOpt) getMulticastInterface(c *socket.Conn) (*net.Interface, error) {
	return nil, errNotImplemented
}

func (so *sockOpt) setMulticastInterface(c *socket.Conn, ifi *net.Interface) error {
	return errNotImplemented
}

func (so *sockOpt) getICMPFilter(c *socket.Conn) (*ICMPFilter, error) {
	return nil, errNotImplemented
}

func (so *sockOpt) setICMPFilter(c *socket.Conn, f *ICMPFilter) error {
	return errNotImplemented
}

func (so *sockOpt) setGroup(c *socket.Conn, ifi *net.Interface, grp net.IP) error {
	return errNotImplemented
}

func (so *sockOpt) setSourceGroup(c *socket.Conn, ifi *net.Interface, grp, src net.IP) error {
	return errNotImplemented
}

func (so *sockOpt) setBPF(c *socket.Conn, f []bpf.RawInstruction) error {
	return errNotImplemented
}
