// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	"os"

	"github.com/xmidt-org/rbus-go"
	"github.com/xmidt-org/rbus-go/client"
)

func run(args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide unix://file or tcp://host:port")
		return nil
	}

	fmt.Println("Creating a new client")
	c, err := client.New(client.Config{
		URL: args[0],
	})
	if err != nil {
		return err
	}

	fmt.Println("Connecting")
	err = c.Connect()
	if err != nil {
		return err
	}
	fmt.Println("Connected")

	fmt.Println("Sending a message")
	err = c.Send(&rbus.Message{
		SeqNum: 55,
		Flags:  rbus.FLAGS_REQUEST,
		Topic:  "_RTROUTED.INBOX.SUBSCRIBE",
	})
	if err != nil {
		return err
	}
	fmt.Println("Message sent")

	fmt.Println("Reading a message")
	got, err := c.Read()
	if err != nil {
		return err
	}
	fmt.Println("Message read")
	fmt.Println(got)

	return nil
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		panic(err)
	}
}
