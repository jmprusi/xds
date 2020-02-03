// Copyright 2018 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package main contains the test driver for testing xDS manually.
package main

import (
	"context"
	"flag"
	"time"

	"github.com/miekg/xds/pkg/cache2"
	"github.com/miekg/xds/pkg/log"
	"github.com/miekg/xds/pkg/server"
)

var (
	nodeID = flag.String("nodeID", "test-id", "Node ID")
	addr   = flag.String("addr", ":18000", "management server address")
	conf   = flag.String("conf", ".", "cluster configuration directory")
)

// main returns code 1 if any of the batches failed to pass all requests
func main() {
	flag.Parse()
	clusters, err := parseClusters(*conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Parsed %d clusters from directory %q", len(clusters), *conf)

	// create a cache
	config := cache2.New()
	for _, cla := range clusters {
		config.Insert(cla)
	}
	log.Info("Initialized cache with 'v1' of cluster info")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	srv := server.NewServer(ctx, config)
	go RunManagementServer(ctx, srv, *addr) // start the xDS server

	for {
		log.Info("Still alive")
		time.Sleep(5 * time.Second)
	}
	// ^c handling: TODO
	cancel()
}
