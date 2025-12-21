package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewResource builds and returns a new OpenTelemetry Resource.
// It merges default, environment, and static attributes.
func NewResource(ctx context.Context, opts ...resource.Option) (*resource.Resource, error) {
	// Add standard resource options plus any user-provided ones.
	options := []resource.Option{
		resource.WithFromEnv(), // reads OTEL_RESOURCE_ATTRIBUTES, etc.
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(getEnv("SERVICE_NAME", "api")),
			semconv.ServiceVersionKey.String(getEnv("SERVICE_VERSION", "1.0.0")),
			semconv.DeploymentEnvironmentKey.String(getEnv("DEPLOYMENT_ENV", "dev")),
		),
	}

	// Append any additional custom options provided by caller.
	options = append(options, opts...)

	res, err := resource.New(ctx, options...)
	return res, err
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
