package main

import "net"

type Client struct {
	conn          net.Conn
	authenticated bool
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}