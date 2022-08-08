/*
Copyright 2019 The Kubernetes Authors.

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

package openstack

import (
	"fmt"
	"net"
)

const (
	openstackExternalIPType = "OS-EXT-IPS:type"
	openstackAddressFixed   = "fixed"
	openstackAddress        = "addr"
)

func getAllServerFixedIPs(addrs map[string]interface{}) []string {
	var fixedIPs []string
	for _, address := range addrs {
		if addresses, ok := address.([]interface{}); ok {
			for _, addr := range addresses {
				addrMap := addr.(map[string]interface{})
				if addrType, ok := addrMap[openstackExternalIPType]; ok && addrType == openstackAddressFixed {
					if fixedIP, ok := addrMap[openstackAddress]; ok {
						if fixedIPStr, ok := fixedIP.(string); ok {
							fixedIPs = append(fixedIPs, fixedIPStr)
						}
					}
				}
			}
		}
	}
	return fixedIPs
}

func GetServerFixedIP(addrs map[string]interface{}, name string, networkCIDRs []*net.IPNet) (poolAddress string, err error) {
	fixedIPs := getAllServerFixedIPs(addrs)

	if networkCIDRs != nil {
		for _, cidr := range networkCIDRs {
			for _, fixedIP := range fixedIPs {
				if cidr.Contains(net.ParseIP(fixedIP)) {
					return fixedIP, nil
				}
			}
		}
	} else if len(fixedIPs) > 0 {
		return fixedIPs[0], nil
	}

	return "", fmt.Errorf("failed to find Fixed IP address for server %s", name)
}
