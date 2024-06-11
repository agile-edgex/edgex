/*******************************************************************************
 * Copyright 2018 Dell Inc.
 * Copyright 2023 Intel Corporation
 * Copyright 2023 IOTech Ltd.
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

package metadata

import (
	"context"
	"os"

	"github.com/agile-edgex/edgex-go"
	"github.com/agile-edgex/edgex-go/internal/core/metadata/config"
	"github.com/agile-edgex/edgex-go/internal/core/metadata/container"
	"github.com/agile-edgex/edgex-go/internal/core/metadata/uom"
	pkgHandlers "github.com/agile-edgex/edgex-go/internal/pkg/bootstrap/handlers"

	"github.com/agile-edgex/go-mod-core-contracts/v3/common"

	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/flags"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/handlers"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/interfaces"
	"github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/startup"
	bootstrapConfig "github.com/agile-edgex/go-mod-bootstrap/v3/config"
	"github.com/agile-edgex/go-mod-bootstrap/v3/di"

	"github.com/labstack/echo/v4"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *echo.Echo) {
	startupTimer := startup.NewStartUpTimer(common.CoreMetaDataServiceKey)

	// All common command-line flags have been moved to DefaultCommonFlags. Service specific flags can be add here,
	// by inserting service specific flag prior to call to commonFlags.Parse().
	// Example:
	// 		flags.FlagSet.StringVar(&myvar, "m", "", "Specify a ....")
	//      ....
	//      flags.Parse(os.Args[1:])
	//
	f := flags.New()
	f.Parse(os.Args[1:])

	configuration := &config.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	httpServer := handlers.NewHttpServer(router, true)

	bootstrap.Run(
		ctx,
		cancel,
		f,
		common.CoreMetaDataServiceKey,
		common.ConfigStemCore,
		configuration,
		startupTimer,
		dic,
		true,
		bootstrapConfig.ServiceTypeOther,
		[]interfaces.BootstrapHandler{
			uom.BootstrapHandler,
			pkgHandlers.NewDatabase(httpServer, configuration, container.DBClientInterfaceName).BootstrapHandler, // add db client bootstrap handler
			handlers.MessagingBootstrapHandler,
			handlers.NewServiceMetrics(common.CoreMetaDataServiceKey).BootstrapHandler, // Must be after Messaging
			NewBootstrap(router, common.CoreMetaDataServiceKey).BootstrapHandler,
			httpServer.BootstrapHandler,
			handlers.NewStartMessage(common.CoreMetaDataServiceKey, edgex.Version).BootstrapHandler,
		},
	)
}
