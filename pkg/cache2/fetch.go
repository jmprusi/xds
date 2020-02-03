package cache2

import (
	"fmt"
	"sort"
	"strconv"

	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/miekg/xds/pkg/resource"
)

// Fetch fetches cluster data from the cluster. Here we probably deviate from the spec, as empty versions are allowed and we
// will return the full list again. For versioning we use the highest version we see in the cache and use that as the version
// in the reply.
func (c *Cluster) Fetch(req *discoverypb.DiscoveryRequest) (*discoverypb.DiscoveryResponse, error) {
	var resources []*any.Any

	switch req.TypeUrl {
	case EndpointType:

	case ClusterType:
		// As we only store ClusterLoadAssignments, we need to create a cluster response.
		sort.Strings(req.ResourceNames)
		clusters := req.ResourceNames
		if len(req.ResourceNames) == 0 {
			clusters = c.All()
		}
		version := uint64(0)
		for _, n := range clusters {
			cla, v := c.Retrieve(n)
			if v > version {
				version = v
			}
			cluster := resource.MakeCluster(cla.GetClusterName())
			data, err := MarshalResource(cluster)
			if err != nil {
				return nil, err
			}
			resources = append(resources, &any.Any{TypeUrl: ClusterType, Value: data})
		}
		versionInfo := strconv.FormatUint(version, 10)
		if versionInfo == req.VersionInfo { // client is up to date
			return nil, SkipFetchError{}
		}
		return &discoverypb.DiscoveryResponse{VersionInfo: versionInfo, Resources: resources, TypeUrl: ClusterType}, nil
	}
	return nil, fmt.Errorf("unrecognized/unsupported type %q:", req.TypeUrl)
}

// SkipFetchError is the error returned when the cache fetch is short
// circuited due to the client's version already being up-to-date.
type SkipFetchError struct{}

// Error satisfies the error interface
func (e SkipFetchError) Error() string {
	return "skip fetch: version up to date"
}
