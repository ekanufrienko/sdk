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

// Package dnscontext provides a dns context specific chain element. It just adds dns configs into dns context of connection context.
// It also provides a possibility to use custom dns configs getters for setup users endpoints.
package dnscontext

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
)

// GetDNSConfigsFunc gets dns configs
type GetDNSConfigsFunc func() []*networkservice.DNSConfig

type dnsContextServer struct {
	getDNSConfigs GetDNSConfigsFunc
}

// NewServer creates dns context chain server element.
func NewServer(getter GetDNSConfigsFunc) networkservice.NetworkServiceServer {
	return &dnsContextServer{getDNSConfigs: getter}
}

func (d *dnsContextServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	if d.getDNSConfigs != nil {
		request.GetConnection().GetContext().DnsContext = &networkservice.DNSContext{
			Configs: d.getDNSConfigs(),
		}
	}
	return next.Server(ctx).Request(ctx, request)
}

func (d *dnsContextServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	if d.getDNSConfigs != nil {
		conn.GetContext().DnsContext = &networkservice.DNSContext{
			Configs: d.getDNSConfigs(),
		}
	}
	return next.Server(ctx).Close(ctx, conn)
}
