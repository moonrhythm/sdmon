package sdmon

import (
	"net/http"
	"strconv"

	"contrib.go.opencensus.io/exporter/stackdriver"
	sdpropagation "go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// Trace returns trace middleware with given sampling probability (0-1)
func Trace(samplingProbability float64) func(h http.Handler) http.Handler {
	if !inited {
		return noopMiddleware
	}

	envSamplingProbability := cfg.String(envTraceSampingProbability)
	if envSamplingProbability != "" {
		samplingProbability, _ = strconv.ParseFloat(envSamplingProbability, 64)
	}

	if samplingProbability <= 0 {
		return noopMiddleware
	}

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:               projectID,
		MonitoringClientOptions: opts,
		TraceClientOptions:      opts,
	})
	if err != nil {
		return noopMiddleware
	}

	trace.RegisterExporter(exporter)

	return func(h http.Handler) http.Handler {
		return &ochttp.Handler{
			Handler:     h,
			Propagation: &sdpropagation.HTTPFormat{},
			FormatSpanName: func(r *http.Request) string {
				proto := r.Header.Get("X-Forwarded-Proto") // trust proxy
				if proto == "" {
					if r.TLS == nil {
						proto = "http"
					} else {
						proto = "https"
					}
				}
				return proto + "://" + r.Host + r.RequestURI
			},
			StartOptions: trace.StartOptions{
				SpanKind: trace.SpanKindServer,
				Sampler:  trace.ProbabilitySampler(samplingProbability),
			},
			IsPublicEndpoint: true,
		}
	}
}

func noopMiddleware(h http.Handler) http.Handler {
	return h
}
