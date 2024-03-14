package intervalLogStats

import (
	"context"
	"errors"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
)

// schema for connector
type logStat struct {
	component.Component
	consumer.Logs

	config          Config
	metricsConsumer consumer.Metrics
	logger          *zap.Logger

	// Include these parameters if a specific implementation for the Start and Shutdown function are not needed
	StartFunc    *component.StartFunc
	ShutdownFunc *component.ShutdownFunc

	metricsMetricDefs    map[string]metricDef[ottlmetric.TransformContext]
	dataPointsMetricDefs map[string]metricDef[ottldatapoint.TransformContext]
	logsMetricDefs       map[string]metricDef[ottllog.TransformContext]
}

// newConnector is a function to create a new connector
func newConnector(logger *zap.Logger, config component.Config) (*logStat, error) {
	logger.Info("Building intervalLogStats connector")
	cfg := config.(*Config)

	return &logStat{
		config: *cfg,
		logger: logger,
	}, nil
}

// Capabilities implements the consumer interface.
func (c *logStat) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs method is called for each instance of a log sent to the connector
func (c *logStat) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	var multiError error
	countMetrics := pmetric.NewMetrics()
	countMetrics.ResourceMetrics().EnsureCapacity(ld.ResourceLogs().Len())
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		counter := newCounter[ottllog.TransformContext](c.logsMetricDefs)

		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLogs := resourceLog.ScopeLogs().At(j)

			for k := 0; k < scopeLogs.LogRecords().Len(); k++ {
				logRecord := scopeLogs.LogRecords().At(k)

				lCtx := ottllog.NewTransformContext(logRecord, scopeLogs.Scope(), resourceLog.Resource())
				multiError = errors.Join(multiError, counter.update(ctx, logRecord.Attributes(), lCtx))
			}
		}

		if len(counter.counts) == 0 {
			continue // don't add an empty resource
		}

		countResource := countMetrics.ResourceMetrics().AppendEmpty()
		resourceLog.Resource().Attributes().CopyTo(countResource.Resource().Attributes())

		countResource.ScopeMetrics().EnsureCapacity(resourceLog.ScopeLogs().Len())
		countScope := countResource.ScopeMetrics().AppendEmpty()
		countScope.Scope().SetName(scopeName)

		counter.appendMetricsTo(countScope.Metrics())
	}
	if multiError != nil {
		return multiError
	}
	return c.metricsConsumer.ConsumeMetrics(ctx, countMetrics)
}
