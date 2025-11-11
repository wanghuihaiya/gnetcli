//go:build !windows

package ssh

import (
	"net"
	"os"
	"syscall"
)

func (m *SSHTunnel) makeSocketFromSocketPair() (net.Conn, net.Conn, error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, err
	}

	f0 := os.NewFile(uintptr(fds[0]), "socketpair-0")
	defer f0.Close()
	c0, err := net.FileConn(f0)
	if err != nil {
		return nil, nil, err
	}
	f1 := os.NewFile(uintptr(fds[1]), "socketpair-0")
	defer f1.Close()
	c1, err := net.FileConn(f1)
	if err != nil {
		return nil, nil, err
	}

	return c0, c1, nil
}
