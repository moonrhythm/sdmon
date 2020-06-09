package sdmon

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

// Report reports error (non-blocking)
// request, user, stack can be empty
func Report(err interface{}, r *http.Request, user string, stack []byte) {
	if errorClient == nil {
		return
	}

	var e error
	if err, ok := err.(error); ok {
		e = err
	} else {
		e = fmt.Errorf("%v", err)
	}

	errorClient.Report(errorreporting.Entry{
		Error: e,
		Req:   r,
		User:  user,
		Stack: stack,
	})
}

// Reportf reports error (non-blocking) using fmt.Errorf
func Reportf(format string, v ...interface{}) {
	Report(fmt.Errorf(format, v...), nil, "", nil)
}
