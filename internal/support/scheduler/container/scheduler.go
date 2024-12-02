//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/agile-edge/edgex/internal/support/scheduler/infrastructure/interfaces"

	"github.com/agile-edge/go-mod-bootstrap/v3/di"
)

// SchedulerManagerName contains the name of the interfaces.SchedulerManager implementation in the DIC.
var SchedulerManagerName = di.TypeInstanceToName((*interfaces.SchedulerManager)(nil))

// SchedulerManagerFrom helper function queries the DIC and returns the interfaces.SchedulerManager implementation.
func SchedulerManagerFrom(get di.Get) interfaces.SchedulerManager {
	return get(SchedulerManagerName).(interfaces.SchedulerManager)
}
