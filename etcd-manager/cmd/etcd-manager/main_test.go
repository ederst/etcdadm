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
	o := getTestData("192.168.0.0/16", "aws")

	expectedErr := fmt.Errorf("is only supported with provider 'openstack'")

	_, actualErr := parseNetworkCIDR(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseNetworkCIDRReturnsErrorOnInvalidCIDR(t *testing.T) {
	invalidCIDR := "192.168.0.0/123"
	o := getTestData(invalidCIDR, "openstack")

	expectedErr := &net.ParseError{Type: "CIDR address", Text: invalidCIDR}

	_, actualErr := parseNetworkCIDR(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseNetworkCIDRReturnsParsedCIDR(t *testing.T) {
	cidr := "192.168.0.0/16"
	o := getTestData(cidr, "openstack")

	_, expectedNetworkCIDR, _ := net.ParseCIDR(cidr)

	actualNetworkCIDR, err := parseNetworkCIDR(o)

	assertTestResults(t, err, expectedNetworkCIDR, actualNetworkCIDR)
}
