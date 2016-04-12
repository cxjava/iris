package iris

import (
	"fmt"
)

var (
	// Server
	ErrServerPortAlreadyUsed   = NewError("Server can't run, port is already used")
	ErrServerAlreadyStarted    = NewError("Server is already started and listening")
	ErrServerOptionsMissing    = NewError("You have to pass iris.ServerOptions")
	ErrServerTlsOptionsMissing = NewError("You have to set CertFile and KeyFile to iris.ServerOptions before ListenTLS")
	ErrServerIsClosed          = NewError("Can't close the server, propably is already closed or never started")
	ErrServerUnknown           = NewError("Unknown error from Server")
	ErrParsedAddr              = NewError("ListeningAddr error, for TCP and UDP, the syntax of ListeningAddr is host:port, like 127.0.0.1:8080. If host is omitted, as in :8080, Listen listens on all available interfaces instead of just the interface with the given host address. See Dial for more details about address syntax")

	// Template
	ErrTemplateParse    = NewError("Couldn't load templates %s")
	ErrTemplateWatch    = NewError("Templates watcher couldn't be started, error: %s")
	ErrTemplateWatching = NewError("Error when watching templates :%s")
)

func NewError(format string, args ...interface{}) error {
	return fmt.Errorf(LoggerIrisPrefix+"Error: "+format, args...)
}

func Printf(logger *Logger, err error, args ...interface{}) {
	if logger.IsEnabled() {
		logger.Printf(err.Error(), args...)
	}

}
