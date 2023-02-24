package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)

	if err != nil {
		log.Fatalf("%s: %v", "failed to create metric exporter", err)
	}

	reader := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(2*time.Second),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(reader),
	)
	global.SetMeterProvider(meterProvider)

	meter := global.Meter("my-meter")

	histogram, err := meter.Int64Histogram("my-histogram")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}

	for {
		time.Sleep(time.Duration(200) * time.Millisecond)
		histogram.Record(ctx, int64(200))
	}
}
