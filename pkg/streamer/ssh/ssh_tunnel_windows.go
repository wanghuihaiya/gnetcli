//go:build windows

package ssh

import (
	"fmt"
	"net"
)

func (m *SSHTunnel) makeSocketFromSocketPair() (net.Conn, net.Conn, error) {
	c1, c2 := net.Pipe()
	if c1 == nil || c2 == nil {
		return nil, nil, fmt.Errorf("failed to create pipe sockets")
	}
	return c1, c2, nil
}
