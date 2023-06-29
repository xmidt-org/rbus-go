// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package client

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/xmidt-org/rbus-go"
)

var (
	ErrInvalidState = errors.New("invalid state")
	ErrInvalidInput = errors.New("invalid input")
)

type Config struct {
	URL string
}

type Client struct {
	network string
	address string
	con     net.Conn
	m       sync.Mutex
}

func New(c Config) (*Client, error) {
	parts := strings.Split(c.URL, "://")
	switch parts[0] {
	case "unix":
		parts[1] = "/" + parts[1]
	case "tcp":
	default:
		return nil, ErrInvalidInput
	}

	return &Client{
		network: parts[0],
		address: parts[1],
	}, nil
}

func (c *Client) Connect() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.con != nil {
		return nil
	}

	con, err := net.Dial(c.network, c.address)
	if err != nil {
		return err
	}

	c.con = con
	return nil
}

func (c *Client) Disconnect() error {
	c.m.Lock()
	defer c.m.Unlock()

	err := c.con.Close()
	c.con = nil
	return err
}

func (c *Client) Send(m *rbus.Message) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.con == nil {
		return ErrInvalidState
	}

	buf, err := m.Encode()
	if err != nil {
		return err
	}

	n, err := c.con.Write(buf)
	if err != nil {
		return err
	}

	if n < len(buf) {
		return errors.New("not all bytes were sent.")
	}

	return nil
}

func (c *Client) Read() (*rbus.Message, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.con == nil {
		return nil, ErrInvalidState
	}

	buf := make([]byte, 4096)

	n, err := c.con.Read(buf)
	if err != nil {
		return nil, err
	}

	buf = buf[0:n]

	return rbus.Decode(buf)
}
