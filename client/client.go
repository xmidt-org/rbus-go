// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package client

import "sync"

type Config struct {
	URL string
}

type Client struct {
	url string
	m   sync.Mutex
}

func New(c Config) *Client {
	return &Client{
		url: c.URL,
	}
}

func (c *Client) Connect() error {
	return nil
}

func (c *Client) Disconnect() error {
	return nil
}

func (c *Client) Send(s string) error {
	c.m.Lock()
	defer c.m.Unlock()
	return nil
}
