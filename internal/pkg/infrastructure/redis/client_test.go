//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	dataInterfaces "github.com/agile-edge/edgex/internal/core/data/infrastructure/interfaces"
	metadataInterfaces "github.com/agile-edge/edgex/internal/core/metadata/infrastructure/interfaces"

	schedulerInterfaces "github.com/agile-edge/edgex/internal/support/scheduler/infrastructure/interfaces"

	notificationsInterfaces "github.com/agile-edge/edgex/internal/support/notifications/infrastructure/interfaces"
)

// Check the implementation of Redis satisfies the DB client
var _ dataInterfaces.DBClient = &Client{}
var _ metadataInterfaces.DBClient = &Client{}
var _ schedulerInterfaces.DBClient = &Client{}
var _ notificationsInterfaces.DBClient = &Client{}
