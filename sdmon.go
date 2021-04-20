// Package sdmon is the singleton for internal monitor to cloud operation
package sdmon

import (
	"context"
	"os"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"github.com/acoshift/configfile"
	"google.golang.org/api/option"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

const (
	defaultProjectID = "moonrhythm-monitor"
)

// env name
const (
	envAutoEnable              = "MOONRHYTHM_SDMON_AUTO_ENABLE"
	envEnable                  = "MOONRHYTHM_SDMON_ENABLE"
	envProjectID               = "MOONRHYTHM_SDMON_PROJECT_ID"
	envService                 = "MOONRHYTHM_SDMON_SERVICE"
	envServiceAccountJSON      = "MOONRHYTHM_SDMON_SERVICE_ACCOUNT_JSON"
	envEnableProfiler          = "MOONRHYTHM_SDMON_ENABLE_PROFILER"
	envTraceSampingProbability = "MOONRHYTHM_SDMON_TRACE_SAMPING_PROBABILITY"
)

var (
	inited      bool
	projectID   string
	errorClient *errorreporting.Client
	logClient   *logging.Client
	logWriter   *logging.Logger
	opts        []option.ClientOption
)

var cfg = configfile.NewEnvReader()

func init() {
	if !cfg.Bool(envAutoEnable) {
		return
	}

	os.Setenv(envEnable, "true")
	Init("", "", "")
}

// Init inits stack driver monitor
func Init(googleProjectID, serviceName, googleServiceAccountJSON string) {
	if inited {
		return
	}

	if !cfg.Bool(envEnable) {
		return
	}

	projectID = cfg.StringDefault(envProjectID, googleProjectID)
	if projectID == "" {
		projectID = defaultProjectID
	}

	serviceName = cfg.StringDefault(envService, serviceName)
	if serviceName == "" {
		serviceName = "service"
	}

	googleServiceAccountJSON = cfg.StringDefault(envServiceAccountJSON, googleServiceAccountJSON)
	if googleServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(googleServiceAccountJSON)))
	}

	ctx := context.Background()

	// profiler is resource expensive
	if cfg.Bool(envEnableProfiler) {
		profiler.Start(profiler.Config{
			Service:   serviceName,
			ProjectID: projectID,
		}, opts...)
	}

	errorClient, _ = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
	}, opts...)

	logClient, _ = logging.NewClient(ctx, "projects/"+projectID, opts...)
	if logClient != nil {
		logWriter = logClient.Logger(serviceName, logging.CommonResource(&mrpb.MonitoredResource{
			Type: "global",
			Labels: map[string]string{
				"project_id": projectID,
			},
		}))
	}

	inited = true
}

// Close closes clients
func Close() {
	if errorClient != nil {
		errorClient.Close()
	}
	if logClient != nil {
		logClient.Close()
	}
}
