package service

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/deduplicator"
	dataStore "github.com/tidepool-org/platform/data/store"
	metricservicesClient "github.com/tidepool-org/platform/metricservices/client"
	"github.com/tidepool-org/platform/service"
	taskStore "github.com/tidepool-org/platform/task/store"
	userservicesClient "github.com/tidepool-org/platform/userservices/client"
)

type Context interface {
	service.Context

	MetricServicesClient() metricservicesClient.Client
	UserServicesClient() userservicesClient.Client

	DataFactory() data.Factory
	DataDeduplicatorFactory() deduplicator.Factory

	DataStoreSession() dataStore.Session
	TaskStoreSession() taskStore.Session

	AuthenticationDetails() userservicesClient.AuthenticationDetails
	SetAuthenticationDetails(authenticationDetails userservicesClient.AuthenticationDetails)
}

type HandlerFunc func(context Context)
