package grpclogger

import (
	"io/ioutil"
	"os"
	"strconv"

	"google.golang.org/grpc/grpclog"
)

func New(severityLevel, verbosityLevel string) grpclog.LoggerV2 {
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	switch severityLevel {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	if vl, err := strconv.Atoi(verbosityLevel); err == nil {
		v = vl
	}
	return grpclog.NewLoggerV2WithVerbosity(infoW, warningW, errorW, v)
}
