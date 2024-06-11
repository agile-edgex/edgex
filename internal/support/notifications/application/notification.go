//
// Copyright (C) 2021-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agile-edgex/edgex-go/internal/pkg/correlation"
	"github.com/agile-edgex/edgex-go/internal/support/notifications/container"

	bootstrapContainer "github.com/agile-edgex/go-mod-bootstrap/v3/bootstrap/container"
	"github.com/agile-edgex/go-mod-bootstrap/v3/di"
	"github.com/agile-edgex/go-mod-core-contracts/v3/dtos"
	"github.com/agile-edgex/go-mod-core-contracts/v3/errors"
	"github.com/agile-edgex/go-mod-core-contracts/v3/models"

	"github.com/google/uuid"
)

var asyncPurgeNotificationOnce sync.Once

// The AddNotification function accepts the new Notification model from the controller function
// and then invokes AddNotification function of infrastructure layer to add new Notification
func AddNotification(n models.Notification, ctx context.Context, dic *di.Container) (id string, edgeXerr errors.EdgeX) {
	dbClient := container.DBClientFrom(dic.Get)
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	addedNotification, edgeXerr := dbClient.AddNotification(n)
	if edgeXerr != nil {
		return "", errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	lc.Debugf("Notification created on DB successfully. Notification ID: %s, Correlation-ID: %s ",
		addedNotification.Id,
		correlation.FromContext(ctx))

	go distribute(dic, addedNotification) // nolint:errcheck

	return addedNotification.Id, nil
}

// NotificationsByCategory queries notifications with offset, limit, and category
func NotificationsByCategory(offset, limit int, category string, dic *di.Container) (notifications []dtos.Notification, totalCount uint32, err errors.EdgeX) {
	if category == "" {
		return notifications, totalCount, errors.NewCommonEdgeX(errors.KindContractInvalid, "category is empty", nil)
	}
	dbClient := container.DBClientFrom(dic.Get)
	notificationModels, err := dbClient.NotificationsByCategory(offset, limit, category)
	if err == nil {
		totalCount, err = dbClient.NotificationCountByCategory(category)
	}
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notifications = make([]dtos.Notification, len(notificationModels))
	for i, n := range notificationModels {
		notifications[i] = dtos.FromNotificationModelToDTO(n)
	}
	return notifications, totalCount, nil
}

// NotificationsByLabel queries notifications with offset, limit, and label
func NotificationsByLabel(offset, limit int, label string, dic *di.Container) (notifications []dtos.Notification, totalCount uint32, err errors.EdgeX) {
	if label == "" {
		return notifications, totalCount, errors.NewCommonEdgeX(errors.KindContractInvalid, "label is empty", nil)
	}
	dbClient := container.DBClientFrom(dic.Get)
	notificationModels, err := dbClient.NotificationsByLabel(offset, limit, label)
	if err == nil {
		totalCount, err = dbClient.NotificationCountByLabel(label)
	}
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notifications = make([]dtos.Notification, len(notificationModels))
	for i, n := range notificationModels {
		notifications[i] = dtos.FromNotificationModelToDTO(n)
	}
	return notifications, totalCount, nil
}

// NotificationById queries notification by ID
func NotificationById(id string, dic *di.Container) (notification dtos.Notification, edgeXerr errors.EdgeX) {
	if id == "" {
		return notification, errors.NewCommonEdgeX(errors.KindContractInvalid, "ID is empty", nil)
	}
	if _, err := uuid.Parse(id); err != nil {
		return notification, errors.NewCommonEdgeX(errors.KindContractInvalid, "ID is not a valid UUID", err)
	}

	dbClient := container.DBClientFrom(dic.Get)
	notificationModel, edgeXerr := dbClient.NotificationById(id)
	if edgeXerr != nil {
		return notification, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	notification = dtos.FromNotificationModelToDTO(notificationModel)
	return notification, nil
}

// NotificationsByStatus queries notifications with offset, limit, and status
func NotificationsByStatus(offset, limit int, status string, dic *di.Container) (notifications []dtos.Notification, totalCount uint32, err errors.EdgeX) {
	if status == "" {
		return notifications, totalCount, errors.NewCommonEdgeX(errors.KindContractInvalid, "status is empty", nil)
	}
	dbClient := container.DBClientFrom(dic.Get)
	notificationModels, err := dbClient.NotificationsByStatus(offset, limit, status)
	if err == nil {
		totalCount, err = dbClient.NotificationCountByStatus(status)
	}
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notifications = make([]dtos.Notification, len(notificationModels))
	for i, n := range notificationModels {
		notifications[i] = dtos.FromNotificationModelToDTO(n)
	}
	return notifications, totalCount, nil
}

// NotificationsByTimeRange query notifications with offset, limit and time range
func NotificationsByTimeRange(start int, end int, offset int, limit int, dic *di.Container) (notifications []dtos.Notification, totalCount uint32, err errors.EdgeX) {
	dbClient := container.DBClientFrom(dic.Get)
	notificationModels, err := dbClient.NotificationsByTimeRange(start, end, offset, limit)
	if err == nil {
		totalCount, err = dbClient.NotificationCountByTimeRange(start, end)
	}
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notifications = make([]dtos.Notification, len(notificationModels))
	for i, n := range notificationModels {
		notifications[i] = dtos.FromNotificationModelToDTO(n)
	}
	return notifications, totalCount, nil
}

// DeleteNotificationById deletes the notification by id and all of its associated transmissions
func DeleteNotificationById(id string, dic *di.Container) errors.EdgeX {
	if id == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "id is empty", nil)
	}
	dbClient := container.DBClientFrom(dic.Get)
	err := dbClient.DeleteNotificationById(id)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// NotificationsBySubscriptionName queries notifications by offset, limit and subscriptionName
func NotificationsBySubscriptionName(offset, limit int, subscriptionName string, dic *di.Container) (notifications []dtos.Notification, totalCount uint32, err errors.EdgeX) {
	if subscriptionName == "" {
		return notifications, totalCount, errors.NewCommonEdgeX(errors.KindContractInvalid, "subscriptionName is empty", nil)
	}
	dbClient := container.DBClientFrom(dic.Get)
	subscription, err := dbClient.SubscriptionByName(subscriptionName)
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notificationModels, err := dbClient.NotificationsByCategoriesAndLabels(offset, limit, subscription.Categories, subscription.Labels)
	if err == nil {
		totalCount, err = dbClient.NotificationCountByCategoriesAndLabels(subscription.Categories, subscription.Labels)
	}
	if err != nil {
		return notifications, totalCount, errors.NewCommonEdgeXWrapper(err)
	}
	notifications = make([]dtos.Notification, len(notificationModels))
	for i, n := range notificationModels {
		notifications[i] = dtos.FromNotificationModelToDTO(n)
	}
	return notifications, totalCount, nil
}

// CleanupNotificationsByAge invokes the infrastructure layer function to remove notifications that are older than age. And the corresponding transmissions will also be deleted
// Age is supposed in milliseconds since modified timestamp.
func CleanupNotificationsByAge(age int64, dic *di.Container) errors.EdgeX {
	dbClient := container.DBClientFrom(dic.Get)

	err := dbClient.CleanupNotificationsByAge(age)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// DeleteProcessedNotificationsByAge invokes the infrastructure layer function to remove processed notifications that are older than age. And the corresponding transmissions will also be deleted
// Age is supposed in milliseconds since modified timestamp.
func DeleteProcessedNotificationsByAge(age int64, dic *di.Container) errors.EdgeX {
	dbClient := container.DBClientFrom(dic.Get)

	err := dbClient.DeleteProcessedNotificationsByAge(age)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// AsyncPurgeNotification purge notifications and related transmissions according to the retention capability.
func AsyncPurgeNotification(interval time.Duration, ctx context.Context, dic *di.Container) {
	asyncPurgeNotificationOnce.Do(func() {
		go func() {
			lc := bootstrapContainer.LoggingClientFrom(dic.Get)
			timer := time.NewTimer(interval)
			for {
				timer.Reset(interval)
				select {
				case <-ctx.Done():
					lc.Info("Exiting notification retention")
					return
				case <-timer.C:
					err := purgeNotification(dic)
					if err != nil {
						lc.Errorf("Failed to purge notifications and transmissions, %v", err)
						break
					}
				}
			}
		}()
	})
}

func purgeNotification(dic *di.Container) errors.EdgeX {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	dbClient := container.DBClientFrom(dic.Get)
	config := container.ConfigurationFrom(dic.Get)
	total, err := dbClient.NotificationTotalCount()
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to query notification total count, %v", err)
	}
	if total >= config.Retention.MaxCap {
		lc.Debugf("Purging the notification amount %d to the minimum capacity %d", total, config.Retention.MinCap)
		// Query the latest notification and clean notifications by modified date.
		notification, err := dbClient.LatestNotificationByOffset(config.Retention.MinCap)
		if err != nil {
			return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("failed to query notification with offset '%d'", config.Retention.MinCap), err)
		}
		now := time.Now().UnixMilli()
		age := now - notification.Modified
		err = dbClient.CleanupNotificationsByAge(age)
		if err != nil {
			return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("failed to delete notifications and transmissions by age '%d'", age), err)
		}
	}
	return nil
}
