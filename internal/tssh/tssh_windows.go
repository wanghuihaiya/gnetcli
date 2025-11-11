//go:build windows

package tssh

import (
	"fmt"
	"net"
)

// makeSocketFromSocketPair is not supported on Windows.
// SSH ControlMaster functionality is Unix/Linux specific.
func makeSocketFromSocketPair() (net.Conn, uintptr, net.Conn, uintptr, error) {
	return nil, 0, nil, 0, fmt.Errorf("socketpair is not supported on Windows")
}

// sendFd is not supported on Windows.
// File descriptor passing is Unix/Linux specific.
func sendFd(conn *net.UnixConn, fd uintptr) error {
	return fmt.Errorf("sendFd is not supported on Windows")
}
