package pkg

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	"net"
	"time"

	"github.com/bborbe/errors"
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

func Ping(ctx context.Context, destination string) error {
	ipAddr, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		return fmt.Errorf("resolve error: %v", err)
	}

	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return fmt.Errorf("dial error (need sudo?): %v", err)
	}
	defer conn.Close()

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

	buffer.Write([]byte("HELLO-PING"))
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], Checksum(b)) // Set checksum

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

	if reply[20] != icmpEchoReply {
		return fmt.Errorf("invalid reply type: got %d, want %d", reply[20], icmpEchoReply)
	}
	glog.V(2).Infof("Reply from %s: bytes=%d time=%v\n", ipAddr.String(), n, duration)
	return nil
}
