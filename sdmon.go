// package sdmon is the singleton for internal monitor to stack driver
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

var (
	inited      bool
	projectID   string
	errorClient *errorreporting.Client
	logClient   *logging.Client
	logWriter   *logging.Logger
	opts        []option.ClientOption
)

var cfg = configfile.NewEnvReader()

// auto init
func init() {
	if !cfg.Bool("MOONRHYTHM_SDMON_AUTO_ENABLE") {
		return
	}

	os.Setenv("MOONRHYTHM_SDMON_ENABLE", "true")
	Init("", "", "")
}

// Init inits stack driver monitor
func Init(googleProjectID, serviceName, googleServiceAccountJSON string) {
	if inited {
		return
	}

	if !cfg.Bool("MOONRHYTHM_SDMON_ENABLE") {
		return
	}

	projectID = cfg.StringDefault("MOONRHYTHM_SDMON_PROJECT_ID", googleProjectID)
	if projectID == "" {
		projectID = defaultProjectID
	}

	serviceName = cfg.StringDefault("MOONRHYTHM_SDMON_SERVICE", serviceName)
	if serviceName == "" {
		serviceName = "service"
	}

	googleServiceAccountJSON = cfg.StringDefault("MOONRHYTHM_SDMON_SERVICE_ACCOUNT_JSON", googleServiceAccountJSON)

	var opts []option.ClientOption
	if googleServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(googleServiceAccountJSON)))
	}

	ctx := context.Background()

	profiler.Start(profiler.Config{
		Service:   serviceName,
		ProjectID: projectID,
	}, opts...)

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
