package intervalLogStats

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

//go:generate mdatagen metadata.yaml

const scope_name: otelcol/intervalLogStats

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/tangentland/intervalLogStats/internal/metadata"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/filterottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
)

const (
	defaultVal = "request.n"
	// this is the name used to refer to the connector in the config.yaml
	typeStr = "intervalLogStats"
)

// NewFactory creates a factory for example connector.
func NewFactory() connector.Factory {
	// OpenTelemetry connector factory to make a factory for connectors

	return connector.NewFactory(
		typeStr,
		createDefaultConfig,
		connector.WithLogsToMetrics(createLogsToMetrics, 3),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		AttributeName: defaultVal,
	}
}

// createLogsToMetrics creates a logs to metrics connector based on provided config.
func createLogsToMetrics(
	ctx context.Context,
	params connector.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (connector.Logs, error) {
	c, err := newConnector(set.Logger, cfg)
	if err != nil {
		return nil, err
	}

	metricDefs := make(map[string]metricDef[ottllog.TransformContext], len(c.Logs))
	for name, info := range c.Logs {
		md := metricDef[ottllog.TransformContext]{
			desc:  info.Description,
			attrs: info.Attributes,
		}
		if len(info.Conditions) > 0 {
			// Error checked in Config.Validate()
			condition, _ := filterottl.NewBoolExprForLog(info.Conditions, filterottl.StandardLogFuncs(), ottl.PropagateError, set.TelemetrySettings)
			md.condition = condition
		}
		metricDefs[name] = md
	}

	return &count{
		metricsConsumer: nextConsumer,
		logsMetricDefs:  metricDefs,
	}, nil
}

type metricDef[K any] struct {
	condition expr.BoolExpr[K]
	desc      string
	attrs     []AttributeConfig
}
