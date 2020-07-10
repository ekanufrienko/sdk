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

package capturecontext

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/registry"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/sdk/pkg/registry/core/next"
)

type contextNSClient struct{}

func (c *contextNSClient) Register(ctx context.Context, in *registry.NetworkService, opts ...grpc.CallOption) (*registry.NetworkService, error) {
	markDoneContext(ctx)
	return next.NetworkServiceRegistryClient(ctx).Register(ctx, in, opts...)
}

func (c *contextNSClient) Find(ctx context.Context, in *registry.NetworkServiceQuery, opts ...grpc.CallOption) (registry.NetworkServiceRegistry_FindClient, error) {
	markDoneContext(ctx)
	return next.NetworkServiceRegistryClient(ctx).Find(ctx, in, opts...)
}

func (c *contextNSClient) Unregister(ctx context.Context, in *registry.NetworkService, opts ...grpc.CallOption) (*empty.Empty, error) {
	markDoneContext(ctx)
	return next.NetworkServiceRegistryClient(ctx).Unregister(ctx, in, opts...)
}

// NewNSRegistryClient - creates a new registry.NetworkServiceRegistryClient chain element that store context
// from the adapter server/client and pass it to the next client/server to avoid the problem with losing
// values from adapted server/client context.
func NewNSRegistryClient() registry.NetworkServiceRegistryClient {
	return &contextNSClient{}
}