package service

import (
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	uuid "github.com/satori/go.uuid"

	"github.com/tidepool-org/platform/log"
)

type LoggerMiddleware struct {
	Logger log.Logger
}

const (
	HTTPHeaderTraceRequest = "X-Tidepool-Trace-Request"
	HTTPHeaderTraceSession = "X-Tidepool-Trace-Session"

	LogTraceRequest = "trace-request"
	LogTraceSession = "trace-session"

	RequestEnvLogger       = "logger"
	RequestEnvTraceRequest = "trace-request"
	RequestEnvTraceSession = "trace-session"

	TraceSessionMaximumLength = 64
)

func (l *LoggerMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(response rest.ResponseWriter, request *rest.Request) {
		traceRequest := strings.Replace(uuid.NewV4().String(), "-", "", -1)
		request.Env[RequestEnvTraceRequest] = traceRequest
		response.Header().Add(HTTPHeaderTraceRequest, traceRequest)
		logger := l.Logger.WithField(LogTraceRequest, traceRequest)

		if traceSession, ok := request.Header[HTTPHeaderTraceSession]; ok {
			if len(traceSession) > TraceSessionMaximumLength {
				traceSession = traceSession[:TraceSessionMaximumLength]
			}
			request.Env[RequestEnvTraceSession] = traceSession
			logger = logger.WithField(LogTraceSession, traceSession)
		}

		request.Env[RequestEnvLogger] = logger

		handler(response, request)
	}
}

func GetRequestLogger(request *rest.Request) log.Logger {
	if request != nil {
		if logger, ok := request.Env[RequestEnvLogger].(log.Logger); ok {
			return logger
		}
	}
	return log.RootLogger()
}

func GetRequestTraceRequest(request *rest.Request) string {
	if request != nil {
		if traceRequest, ok := request.Env[RequestEnvTraceRequest].(string); ok {
			return traceRequest
		}
	}
	return ""
}

func GetRequestTraceSession(request *rest.Request) string {
	if request != nil {
		if traceSession, ok := request.Env[RequestEnvTraceSession].(string); ok {
			return traceSession
		}
	}
	return ""
}