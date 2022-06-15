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

func getTestData(ipFilter string, volumeProviderID string) *EtcdManagerOptions {
	var o EtcdManagerOptions
	o.InitDefaults()

	o.IPFilter = ipFilter
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

func TestParseIPFilterReturnsNilByDefault(t *testing.T) {
	o := getTestData("", "")

	_, actualErr := parseIPFilter(o)

	assertTestResults(t, nil, nil, actualErr)
}

func TestParseIPFilterReturnsUnsupportedProviderError(t *testing.T) {
	o := getTestData("192.168.0.0/16", "")

	expectedErr := fmt.Errorf("is only supported with provider 'openstack'")

	_, actualErr := parseIPFilter(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseIPFilterReturnsErrorOnInvalidCIDR(t *testing.T) {
	o := getTestData("192.168.0.0/123", "openstack")

	expectedErr := &net.ParseError{Type: "CIDR address", Text: o.IPFilter}

	_, actualErr := parseIPFilter(o)

	assertTestResults(t, nil, expectedErr, actualErr)
}

func TestParseIPFilterReturnsParsedCIDR(t *testing.T) {
	o := getTestData("192.168.0.0/16", "openstack")

	_, expectedIPFilter, _ := net.ParseCIDR(o.IPFilter)

	actualIPFilter, err := parseIPFilter(o)

	assertTestResults(t, err, expectedIPFilter, actualIPFilter)
}

func TestParseIPFilterReturnsParsedIPv6CIDR(t *testing.T) {
	o := getTestData("2001:db8::/64", "openstack")

	_, expectedIPFilter, _ := net.ParseCIDR(o.IPFilter)

	actualIPFilter, err := parseIPFilter(o)

	assertTestResults(t, err, expectedIPFilter, actualIPFilter)
}
