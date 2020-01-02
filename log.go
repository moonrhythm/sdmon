package sdmon

import (
	"fmt"

	"cloud.google.com/go/logging"
)

type Severity logging.Severity

const (
	Debug    = Severity(logging.Debug)
	Info     = Severity(logging.Info)
	Warning  = Severity(logging.Warning)
	Error    = Severity(logging.Error)
	Critical = Severity(logging.Critical)
)

// Log logs payload, payload can be string or can marshal into json
func Log(severity Severity, payload interface{}) {
	if logWriter == nil {
		return
	}

	logWriter.Log(logging.Entry{
		Severity: logging.Severity(severity),
		Payload:  payload,
	})
}

// Logf logs string
func Logf(severity Severity, s string, v ...interface{}) {
	Log(severity, fmt.Sprintf(s, v...))
}

// LogDebug logs payload with debug severity
func LogDebug(payload interface{}) {
	Log(Debug, payload)
}

// LogDebugf logs string with debug severity
func LogDebugf(s string, v ...interface{}) {
	Logf(Debug, s, v...)
}

// LogInfo logs payload with info severity
func LogInfo(payload interface{}) {
	Log(Info, payload)
}

// LogInfof logs string with info severity
func LogInfof(s string, v ...interface{}) {
	Logf(Info, s, v...)
}

// LogWarning logs payload with warning severity
func LogWarning(payload interface{}) {
	Log(Warning, payload)
}

// LogWarningf logs string with warning severity
func LogWarningf(s string, v ...interface{}) {
	Logf(Warning, s, v...)
}

// LogError logs payload with error severity
func LogError(payload interface{}) {
	Log(Error, payload)
}

// LogErrorf logs string with error severity
func LogErrorf(s string, v ...interface{}) {
	Logf(Error, s, v...)
}

// LogCritical logs payload with critical severity
func LogCritical(payload interface{}) {
	Log(Critical, payload)
}

// LogCriticalf logs string with critical severity
func LogCriticalf(s string, v ...interface{}) {
	Logf(Critical, s, v...)
}
