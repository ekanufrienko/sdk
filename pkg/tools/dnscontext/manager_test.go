// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dnscontext

import (
	"testing"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/stretchr/testify/require"
)

func TestManager_StoreAnyDomain(t *testing.T) {
	const expected = `. {
	forward . IP1 IP2
	log
	reload
}`
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		DnsServerIps: []string{"IP1", "IP2"},
	})
	require.Equal(t, m.String(), expected)
}

func TestManager_StoreAnyDomainConflict(t *testing.T) {
	const expected = `. {
	fanout . IP1 IP2 IP3
	log
	reload
}`
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		DnsServerIps: []string{"IP1", "IP2"},
	})
	m.Store("1", &networkservice.DNSConfig{
		DnsServerIps: []string{"IP3"},
	})
	actual := m.String()
	require.Len(t, actual, len(expected))
	require.Contains(t, actual, "IP1")
	require.Contains(t, actual, "IP2")
}

func TestManager_Store(t *testing.T) {
	expected := []string{`zone-a {
	forward . IP1 IP2
	log
}`, `zone-b zone-c {
	forward . IP3 IP4
	log
}`}
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP1", "IP2"},
	})
	m.Store("1", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-b", "zone-c"},
		DnsServerIps:  []string{"IP3", "IP4"},
	})
	actual := m.String()
	require.Contains(t, actual, expected[0])
	require.Contains(t, actual, expected[1])
	require.Len(t, actual, len(expected[0])+len(expected[1])+len("\n"))
}

func TestManager_StoreConflict(t *testing.T) {
	const expected = `zone-a {
	fanout . IP1 IP2 IP3
	log
}`
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP1", "IP2"},
	})
	m.Store("1", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP3", "IP1"},
	})
	actual := m.String()
	require.Len(t, actual, len(expected))
	require.Contains(t, actual, "IP1 IP2")
	require.Contains(t, actual, "IP3")
}

func TestManger_Remove(t *testing.T) {
	const expected = `zone-a {
	forward . IP1 IP2
	log
}`
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP1", "IP2"},
	})
	m.Store("1", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-b", "zone-c"},
		DnsServerIps:  []string{"IP3", "IP4"},
	})
	m.Remove("1")
	require.Equal(t, m.String(), expected)
}
func TestManger_RemoveConflict(t *testing.T) {
	const expected = `zone-a {
	forward . IP1 IP2
	log
}`
	m := NewManager()
	m.Store("0", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP1", "IP2"},
	})
	m.Store("1", &networkservice.DNSConfig{
		SearchDomains: []string{"zone-a"},
		DnsServerIps:  []string{"IP3", "IP1"},
	})
	require.NotEqual(t, m.String(), expected)
	m.Remove("1")
	require.Equal(t, m.String(), expected)
}
