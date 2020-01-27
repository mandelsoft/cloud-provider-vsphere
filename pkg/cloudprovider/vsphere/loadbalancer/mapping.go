/*
 Copyright 2020 The Kubernetes Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

package loadbalancer

import (
	"fmt"
	"strconv"

	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	corev1 "k8s.io/api/core/v1"
)

// Mapping defines the port mapping and protocol
type Mapping struct {
	// SourcePort is the service source port
	SourcePort int
	// NodePort is the service node port
	NodePort int
	// Protoocl is the protocol on the service port
	Protocol corev1.Protocol
}

// NewMapping creates a new Mapping for the given service port
func NewMapping(servicePort corev1.ServicePort) Mapping {
	return Mapping{
		SourcePort: int(servicePort.Port),
		NodePort:   int(servicePort.NodePort),
		Protocol:   servicePort.Protocol,
	}
}

func (m Mapping) String() string {
	return fmt.Sprintf("%s/%d->%d", m.Protocol, m.SourcePort, m.NodePort)
}

// MatchVirtualServer returns true if source port and protocol are matching
func (m Mapping) MatchVirtualServer(server *loadbalancer.LbVirtualServer) bool {
	return server.Port == formatPort(m.SourcePort) && server.IpProtocol == string(m.Protocol)
}

// MatchPool returns true if the pool has the correct port tag
func (m Mapping) MatchPool(pool *loadbalancer.LbPool) bool {
	return checkTags(pool.Tags, portTag(m))
}

// MatchTCPMonitor returns true if the monitor has the correct port tag
func (m Mapping) MatchTCPMonitor(monitor *loadbalancer.LbTcpMonitor) bool {
	return checkTags(monitor.Tags, portTag(m))
}

// MatchNodePort returns true if the server pool member port is equal to the mapping's node port
func (m Mapping) MatchNodePort(server *loadbalancer.LbVirtualServer) bool {
	return server.DefaultPoolMemberPort == formatPort(m.NodePort)
}

func formatPort(port int) string {
	return strconv.FormatInt(int64(port), 10)
}
