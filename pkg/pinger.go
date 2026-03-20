// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
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

// Ping sends an ICMP echo request to the given IP address and waits for a reply.
func Ping(ctx context.Context, ipAddr *net.IPAddr) error {
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return errors.Errorf(ctx, "dial error (need sudo?): %v", err)
	}
	defer conn.Close() //nolint:errcheck

	icmp := icmpMessage{
		Type: icmpEchoRequest,
		Code: 0,
		ID:   0x1234,
		Seq:  1,
	}

	var buffer bytes.Buffer

	if err := binary.Write(&buffer, binary.BigEndian, icmp); err != nil {
		return errors.Wrapf(ctx, err, "write icmp message failed")
	}

	buffer.WriteString("HELLO-PING")
	icmpBytes := buffer.Bytes()
	binary.BigEndian.PutUint16(icmpBytes[2:], Checksum(icmpBytes)) // Set checksum

	// Send ICMP Echo Request
	start := time.Now()

	if _, err := conn.Write(icmpBytes); err != nil {
		return errors.Errorf(ctx, "send error: %v", err)
	}

	// Read ICMP Echo Reply
	reply := make([]byte, 1024)

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return errors.Wrapf(ctx, err, "SetReadDeadline failed")
	}

	bytesRead, err := conn.Read(reply)
	if err != nil {
		return errors.Errorf(ctx, "read timeout or error: %v", err)
	}

	duration := time.Since(start)

	if reply[20] != icmpEchoReply {
		return errors.Errorf(ctx, "invalid reply type: got %d, want %d", reply[20], icmpEchoReply)
	}

	glog.V(2).Infof("Reply from %s: bytes=%d time=%.4fms\n",
		ipAddr.String(), bytesRead, duration.Seconds()*1000)

	return nil
}
