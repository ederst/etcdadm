/*
Copyright 2022 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
)

func getTestData(networkCIDR string, volumeProviderID string) *EtcdManagerOptions {
	var o EtcdManagerOptions
	o.InitDefaults()

	o.NetworkCIDR = networkCIDR
	o.VolumeProviderID = volumeProviderID

	return &o
}

func assertTestResults(t *testing.T, err error, expected interface{}, actual interface{}) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected '%+v', but got '%+v'", expected, actual)
	}
}

func TestParseNetworkCIDRReturnsNilByDefault(t *testing.T) {
	o := getTestData("", "")

	_, actualErr := parseNetworkCIDR(o)

	assertTestResults(t, nil, nil, actualErr)
}

func TestParseNetworkCIDRReturnsUnsupportedProviderError(t *testing.T) {
	o := getTestData("192.168.0.0/16", "")

	expectedErr := fmt.Errorf("is only supported with provider 'openstack'")

	_, actualErr := parseNetworkCIDR(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseNetworkCIDRReturnsErrorOnInvalidCIDR(t *testing.T) {
	o := getTestData("192.168.0.0/123, 2001:db8::/64", "openstack")

	expectedErr := &net.ParseError{Type: "CIDR address", Text: "192.168.0.0/123"}

	_, actualErr := parseNetworkCIDR(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseNetworkCIDRReturnsParsedCIDR(t *testing.T) {
	o := getTestData("192.168.0.0/16, 2001:db8::/64", "openstack")

	var expectedNetworkCIDRs []*net.IPNet
	_, cidr1, _ := net.ParseCIDR("192.168.0.0/16")
	_, cidr2, _ := net.ParseCIDR("2001:db8::/64")
	expectedNetworkCIDRs = append(expectedNetworkCIDRs, cidr1, cidr2)

	actualNetworkCIDRs, err := parseNetworkCIDR(o)

	assertTestResults(t, err, expectedNetworkCIDRs, actualNetworkCIDRs)
}

func TestParseInitDefaultReturnsEmptyStringForNetworkCIDRs(t *testing.T) {
	var o EtcdManagerOptions
	o.InitDefaults()

	assertTestResults(t, nil, "", o.NetworkCIDR)
}

func TestParseInitDefaultReturnsValueOfEnvVarForNetworkCIDRs(t *testing.T) {
	expectedNetworkCIDR := "192.168.0.0/16, 2001:db8::/64"
	os.Setenv("ETCD_MANAGER_NETWORK_CIDR", expectedNetworkCIDR)

	var o EtcdManagerOptions
	o.InitDefaults()

	assertTestResults(t, nil, expectedNetworkCIDR, o.NetworkCIDR)
}
