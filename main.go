// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bborbe/ping/pkg"
	"github.com/golang/glog"
	"os"
	"runtime"
	"time"
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
		fmt.Println("ping: missing host operand")
		os.Exit(1)
	}

	for _, arg := range args {
		go func(ctx context.Context, arg string) {
			ch := time.NewTicker(time.Second).C
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					if err := pkg.Ping(ctx, arg); err != nil {
						fmt.Println("Ping failed:", err)
					}
				}
			}
		}(ctx, arg)
	}
	select {
	case <-ctx.Done():
		fmt.Println("shutting down")
		os.Exit(1)
	}
}
