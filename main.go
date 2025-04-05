// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bborbe/errors"
	"github.com/golang/glog"
)

const (
	icmpEchoRequest = 8
	icmpEchoReply   = 0
)

type icmpMessage struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
}

func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func ping(ctx context.Context, destination string) error {
	// Resolve address
	ipAddr, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		return fmt.Errorf("resolve error: %v", err)
	}

	// Open raw socket
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return fmt.Errorf("dial error (need sudo?): %v", err)
	}
	defer conn.Close()

	// Prepare ICMP Echo Request
	icmp := icmpMessage{
		Type: icmpEchoRequest,
		Code: 0,
		ID:   0x1234,
		Seq:  1,
	}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	// Add some data
	buffer.Write([]byte("HELLO-PING"))
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], checksum(b)) // Set checksum

	// Send ICMP Echo Request
	start := time.Now()
	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	// Read ICMP Echo Reply
	reply := make([]byte, 1024)

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return errors.Wrapf(ctx, err, "SetReadDeadline failed")
	}

	n, err := conn.Read(reply)
	if err != nil {
		return fmt.Errorf("read timeout or error: %v", err)
	}
	duration := time.Since(start)

	// Validate reply
	if reply[20] != icmpEchoReply {
		return fmt.Errorf("invalid reply type: got %d, want %d", reply[20], icmpEchoReply)
	}

	fmt.Printf("Reply from %s: bytes=%d time=%v\n", ipAddr.String(), n, duration)
	return nil
}

func main() {
	ctx := contextWithSig(context.Background())

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("ping: missing host operand")
		os.Exit(1)
	}

	ch := time.NewTicker(time.Second).C
	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down")
			os.Exit(1)
		case <-ch:
			if err := ping(ctx, args[0]); err != nil {
				fmt.Println("Ping failed:", err)
			}
		}
	}
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		defer close(signalCh)

		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case signal, ok := <-signalCh:
			if !ok {
				glog.V(2).Infof("signal channel closed => cancel context ")
				return
			}
			glog.V(2).Infof("got signal %s => cancel context ", signal)
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
