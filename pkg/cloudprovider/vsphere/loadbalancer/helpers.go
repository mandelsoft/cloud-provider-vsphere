/*
 * Copyright 2020 The Kubernetes Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 */

package loadbalancer

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// ObjectName is a namespace/name pair
type ObjectName struct {
	// Namespace is the namespace
	Namespace string
	// Name is the name
	Name string
}

func objectNameFromService(service *corev1.Service) ObjectName {
	return ObjectName{Namespace: service.Namespace, Name: service.Name}
}

func (o ObjectName) String() string {
	return o.Namespace + "/" + o.Name
}

func parseObjectName(name string) ObjectName {
	parts := strings.Split(name, "/")
	return ObjectName{Namespace: parts[0], Name: parts[1]}
}

func stringsEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func collectNodeInternalAddresses(nodes []*corev1.Node) map[string]string {
	set := map[string]string{}
	for _, node := range nodes {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				set[addr.Address] = node.Name
				break
			}
		}
	}
	return set
}
