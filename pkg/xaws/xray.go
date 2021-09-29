package xaws

import (
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xraylog"
	log "github.com/julb/go/pkg/logging"
)

type XrayLoggerProxy struct {
	logger *log.LogWithContext
}

func NewXRayLoggerProxy() *XrayLoggerProxy {
	return &XrayLoggerProxy{}
}

func (proxy *XrayLoggerProxy) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	switch level {
	case xraylog.LogLevelDebug:
		proxy.logger.Debug(msg.String())
	case xraylog.LogLevelInfo:
		proxy.logger.Info(msg.String())
	case xraylog.LogLevelWarn:
		proxy.logger.Warn(msg.String())
	case xraylog.LogLevelError:
		proxy.logger.Error(msg.String())
	}
}
