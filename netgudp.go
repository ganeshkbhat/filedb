// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package netgudp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

const (
	defaultServer = "default"
)

// Server is the UDP server.
type Server struct {
	mu      sync.Mutex  // Used for Server.listen concurrent safety. -- The golang test with data race checks this.
	conn    *Conn       // UDP server connection object.
	address string      // UDP server listening address.
	handler func(*Conn) // Handler for UDP connection.
}

var (
	// serverMapping is used for instance name to its UDP server mappings.
	serverMapping = gmap.NewStrAnyMap(true)
)

// GetServer creates and returns an UDP server instance with given name.
func GetServer(name ...interface{}) *Server {
	serverName := defaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Set(serverName, s)
	return s
}

// NewServer creates and returns an UDP server.
// The optional parameter `name` is used to specify its name, which can be used for
// GetServer function to retrieve its instance.
func NewServer(address string, handler func(*Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}

// SetAddress sets the server address for UDP server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// SetHandler sets the connection handler for UDP server.
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// Close closes the connection.
// It will make server shutdowns immediately.
func (s *Server) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err = s.conn.Close()
	if err != nil {
		err = gerror.Wrap(err, "connection failed")
	}
	return
}

// Run starts listening UDP connection.
func (s *Server) Run() error {
	if s.handler == nil {
		err := gerror.NewCode(gcode.CodeMissingConfiguration, "start running failed: socket handler not defined")
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		err = gerror.Wrapf(err, `net.ResolveUDPAddr failed for address "%s"`, s.address)
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		err = gerror.Wrapf(err, `net.ListenUDP failed for address "%s"`, s.address)
		return err
	}
	s.mu.Lock()
	s.conn = NewConnByNetConn(conn)
	s.mu.Unlock()
	s.handler(s.conn)
	return nil
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *Server) GetListenedAddress() string {
	if !gstr.Contains(s.address, FreePortAddress) {
		return s.address
	}
	var (
		address      = s.address
		listenedPort = s.GetListenedPort()
	)
	address = gstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *Server) GetListenedPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ln := s.conn; ln != nil {
		return ln.LocalAddr().(*net.UDPAddr).Port
	}
	return -1
}

// NewNetConn creates and returns a *net.UDPConn with given addresses.
func NewNetConn(remoteAddress string, localAddress ...string) (*net.UDPConn, error) {
	var (
		err        error
		remoteAddr *net.UDPAddr
		localAddr  *net.UDPAddr
		network    = `udp`
	)
	remoteAddr, err = net.ResolveUDPAddr(network, remoteAddress)
	if err != nil {
		return nil, gerror.Wrapf(
			err,
			`net.ResolveUDPAddr failed for network "%s", address "%s"`,
			network, remoteAddress,
		)
	}
	if len(localAddress) > 0 {
		localAddr, err = net.ResolveUDPAddr(network, localAddress[0])
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ResolveUDPAddr failed for network "%s", address "%s"`,
				network, localAddress[0],
			)
		}
	}
	conn, err := net.DialUDP(network, localAddr, remoteAddr)
	if err != nil {
		return nil, gerror.Wrapf(
			err,
			`net.DialUDP failed for network "%s", local "%s", remote "%s"`,
			network, localAddr.String(), remoteAddr.String(),
		)
	}
	return conn, nil
}

// Send writes data to `address` using UDP connection and then closes the connection.
// Note that it is used for short connection usage.
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendRecv writes data to `address` using UDP connection, reads response and then closes the connection.
// Note that it is used for short connection usage.
func SendRecv(address string, data []byte, receive int, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecv(data, receive, retry...)
}

// MustGetFreePort performs as GetFreePort, but it panics if any error occurs.
func MustGetFreePort() (port int) {
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}
	return port
}

// GetFreePort retrieves and returns a port that is free.
func GetFreePort() (port int, err error) {
	var (
		network = `udp`
		address = `:0`
	)
	resolvedAddr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ResolveUDPAddr failed for network "%s", address "%s"`,
			network, address,
		)
	}
	l, err := net.ListenUDP(network, resolvedAddr)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ListenUDP failed for network "%s", address "%s"`,
			network, resolvedAddr.String(),
		)
	}
	port = l.LocalAddr().(*net.UDPAddr).Port
	_ = l.Close()
	return
}

// GetFreePorts retrieves and returns specified number of ports that are free.
func GetFreePorts(count int) (ports []int, err error) {
	var (
		network = `udp`
		address = `:0`
	)
	for i := 0; i < count; i++ {
		resolvedAddr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ResolveUDPAddr failed for network "%s", address "%s"`,
				network, address,
			)
		}
		l, err := net.ListenUDP(network, resolvedAddr)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ListenUDP failed for network "%s", address "%s"`,
				network, resolvedAddr.String(),
			)
		}
		ports = append(ports, l.LocalAddr().(*net.UDPAddr).Port)
		_ = l.Close()
	}
	return ports, nil
}

// Conn handles the UDP connection.
type Conn struct {
	*net.UDPConn                 // Underlying UDP connection.
	remoteAddr     *net.UDPAddr  // Remote address.
	deadlineRecv   time.Time     // Timeout point for reading data.
	deadlineSend   time.Time     // Timeout point for writing data.
	bufferWaitRecv time.Duration // Interval duration for reading buffer.
}

const (
	defaultRetryInterval  = 100 * time.Millisecond // Retry interval.
	defaultReadBufferSize = 1024                   // (Byte)Buffer size.
	receiveAllWaitTimeout = time.Millisecond       // Default interval for reading buffer.
)

type Retry struct {
	Count    int           // Max retry count.
	Interval time.Duration // Retry interval.
}

// NewConn creates UDP connection to `remoteAddress`.
// The optional parameter `localAddress` specifies the local address for connection.
func NewConn(remoteAddress string, localAddress ...string) (*Conn, error) {
	if conn, err := NewNetConn(remoteAddress, localAddress...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnByNetConn creates an UDP connection object with given *net.UDPConn object.
func NewConnByNetConn(udp *net.UDPConn) *Conn {
	return &Conn{
		UDPConn:        udp,
		deadlineRecv:   time.Time{},
		deadlineSend:   time.Time{},
		bufferWaitRecv: receiveAllWaitTimeout,
	}
}

// Send writes data to remote address.
func (c *Conn) Send(data []byte, retry ...Retry) (err error) {
	for {
		if c.remoteAddr != nil {
			_, err = c.WriteToUDP(data, c.remoteAddr)
		} else {
			_, err = c.Write(data)
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				return err
			}
			// Still failed even after retrying.
			if len(retry) == 0 || retry[0].Count == 0 {
				err = gerror.Wrap(err, `Write data failed`)
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv receives and returns data from remote address.
// The parameter `buffer` is used for customizing the receiving buffer size. If `buffer` <= 0,
// it uses the default buffer size, which is 1024 byte.
//
// There's package border in UDP protocol, we can receive a complete package if specified
// buffer size is big enough. VERY NOTE that we should receive the complete package in once
// or else the leftover package data would be dropped.
func (c *Conn) Recv(buffer int, retry ...Retry) ([]byte, error) {
	var (
		err        error        // Reading error
		size       int          // Reading size
		data       []byte       // Buffer object
		remoteAddr *net.UDPAddr // Current remote address for reading
	)
	if buffer > 0 {
		data = make([]byte, buffer)
	} else {
		data = make([]byte, defaultReadBufferSize)
	}
	for {
		size, remoteAddr, err = c.ReadFromUDP(data)
		if err == nil {
			c.remoteAddr = remoteAddr
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			if len(retry) > 0 {
				// It fails even it retried.
				if retry[0].Count == 0 {
					break
				}
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			err = gerror.Wrap(err, `ReadFromUDP failed`)
			break
		}
		break
	}
	return data[:size], err
}

// SendRecv writes data to connection and blocks reading response.
func (c *Conn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err != nil {
		return nil, err
	}
	return c.Recv(receive, retry...)
}

// RecvWithTimeout reads data from remote address with timeout.
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err = c.SetDeadlineRecv(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer func() {
		_ = c.SetDeadlineRecv(time.Time{})
	}()
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to connection with timeout.
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err = c.SetDeadlineSend(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer func() {
		_ = c.SetDeadlineSend(time.Time{})
	}()
	err = c.Send(data, retry...)
	return
}

// SendRecvWithTimeout writes data to connection and reads response with timeout.
func (c *Conn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err != nil {
		return nil, err
	}
	return c.RecvWithTimeout(receive, timeout, retry...)
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (c *Conn) SetDeadline(t time.Time) (err error) {
	if err = c.UDPConn.SetDeadline(t); err == nil {
		c.deadlineRecv = t
		c.deadlineSend = t
	} else {
		err = gerror.Wrapf(err, `SetDeadline for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineRecv sets the read deadline associated with the connection.
func (c *Conn) SetDeadlineRecv(t time.Time) (err error) {
	if err = c.SetReadDeadline(t); err == nil {
		c.deadlineRecv = t
	} else {
		err = gerror.Wrapf(err, `SetDeadlineRecv for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineSend sets the deadline of sending for current connection.
func (c *Conn) SetDeadlineSend(t time.Time) (err error) {
	if err = c.SetWriteDeadline(t); err == nil {
		c.deadlineSend = t
	} else {
		err = gerror.Wrapf(err, `SetDeadlineSend for connection failed with "%s"`, t)
	}
	return err
}

// SetBufferWaitRecv sets the buffer waiting timeout when reading all data from connection.
// The waiting duration cannot be too long which might delay receiving data from remote address.
func (c *Conn) SetBufferWaitRecv(d time.Duration) {
	c.bufferWaitRecv = d
}

// RemoteAddr returns the remote address of current UDP connection.
// Note that it cannot use c.conn.RemoteAddr() as it is nil.
func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}
