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

package main

import (
	"context"
	"net"

	xdspb2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	cdspb "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverypb2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	edspb "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	healthpb "github.com/envoyproxy/go-control-plane/envoy/service/health/v3"
	"github.com/miekg/xds/pkg/log"
	"github.com/miekg/xds/pkg/server"
	"google.golang.org/grpc"
)

const grpcMaxConcurrentStreams = 1000000

// RunManagementServer starts an xDS server at the given port.
func RunManagementServer(ctx context.Context, server server.Server, server2 server.Server2, addr string) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems.
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// register services
	discoverypb.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	edspb.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	cdspb.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	healthpb.RegisterHealthDiscoveryServiceServer(grpcServer, server)

	discoverypb2.RegisterAggregatedDiscoveryServiceServer(grpcServer, server2)
	xdspb2.RegisterClusterDiscoveryServiceServer(grpcServer, server2)

	log.Infof("Management server listening on %s", addr)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}
