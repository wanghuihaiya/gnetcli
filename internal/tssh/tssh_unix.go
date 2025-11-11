//go:build !windows

package tssh

import (
	"net"
	"os"
	"syscall"
)

func makeSocketFromSocketPair() (net.Conn, uintptr, net.Conn, uintptr, error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, 0, nil, 0, err
	}

	f0 := os.NewFile(uintptr(fds[0]), "socketpair-0")
	c0, err := net.FileConn(f0)
	if err != nil {
		f0.Close()
		return nil, 0, nil, 0, err
	}
	f1 := os.NewFile(uintptr(fds[1]), "socketpair-0")
	c1, err := net.FileConn(f1)
	if err != nil {
		f0.Close()
		f1.Close()
		return nil, 0, nil, 0, err
	}

	return c0, f0.Fd(), c1, f1.Fd(), nil
}

func sendFd(conn *net.UnixConn, fd uintptr) error {
	oob := syscall.UnixRights(int(fd))
	_, _, err := conn.WriteMsgUnix([]byte{}, oob, nil)
	return err
}
