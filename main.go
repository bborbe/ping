// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"

	"github.com/bborbe/ping/pkg"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("v", "2")

	ctx := pkg.ContextWithSig(context.Background())

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("missing host operand")
		os.Exit(1)
	}

	for _, destination := range args {
		ipAddr, err := net.ResolveIPAddr("ip4", destination)
		if err != nil {
			fmt.Printf("resolve error: %v\n", err)
			os.Exit(1)
		}

		go func(ctx context.Context, ipAddr *net.IPAddr) {
			ch := time.NewTicker(time.Second).C
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					if err := pkg.Ping(ctx, ipAddr); err != nil {
						fmt.Println("Ping failed:", err)
					}
				}
			}
		}(ctx, ipAddr)
	}
	select {
	case <-ctx.Done():
		fmt.Println("shutting down")
		os.Exit(1)
	}
}
