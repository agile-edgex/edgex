/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019-2023 Intel Corporation
 * Copyright (C) 2021 IOTech Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package command

import (
	"context"
	"sync"

	"github.com/agile-edgex/edgex/internal/core/command/container"
	bootstrapContainer "github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/container"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/secret"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/startup"
	"github.com/agile-edgex/go-mod-bootstrap/v3/di"
	clients "github.com/agile-edgex/go-mod-core-contracts/v3/clients/http"

	"github.com/labstack/echo/v4"
)

// Bootstrap contains references to dependencies required by the BootstrapHandler.
type Bootstrap struct {
	router      *echo.Echo
	serviceName string
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(router *echo.Echo, serviceName string) *Bootstrap {
	return &Bootstrap{
		router:      router,
		serviceName: serviceName,
	}
}

// BootstrapHandler fulfills the BootstrapHandler contract and performs initialization needed by the command service.
func (b *Bootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	LoadRestRoutes(b.router, dic, b.serviceName)
	config := container.ConfigurationFrom(dic.Get)

	// DeviceServiceCommandClient is not part of the common clients handled by the NewClientsBootstrap handler
	dic.Update(di.ServiceConstructorMap{
		bootstrapContainer.DeviceServiceCommandClientName: func(get di.Get) interface{} { // add API DeviceServiceCommandClient
			jwtSecretProvider := secret.NewJWTSecretProvider(bootstrapContainer.SecretProviderExtFrom(get))
			return clients.NewDeviceServiceCommandClient(jwtSecretProvider, config.Service.EnableNameFieldEscape)
		},
	})

	return true
}
