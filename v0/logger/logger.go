package logger

import (
    "os"
    "io"
    "fmt"
    "log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Logger struct {
    Request *os.File
    Response *os.File

    RequestHandler gin.HandlerFunc
    ResponseHandler gin.HandlerFunc
}

var LoggerInst Logger

func (log *Logger)Init() {

	var err error
	log.Request , err = os.Create(os.Getenv("REQUEST_LOG")+".request.log")
    if err != nil {
        fmt.Println("Error creating log file for requests", err)
    }

	log.Response, err = os.Create(os.Getenv("RESPONSE_LOG")+".response.log")
    if err != nil {
        fmt.Println("Error creating log file for response", err)
    }

	configRequest := sloggin.Config{
		WithRequestBody: true,
		WithResponseBody: false,
		WithRequestHeader: true,
		WithResponseHeader: false,
	}

	configResponse := sloggin.Config{
		WithRequestBody: false,
		WithResponseBody: true,
		WithRequestHeader: false,
		WithResponseHeader: true,
	}

    handlerRequest := slog.New(slog.NewJSONHandler(io.MultiWriter(log.Request), nil))
    handlerResponse := slog.New(slog.NewJSONHandler(io.MultiWriter(log.Response), nil))

    log.RequestHandler = sloggin.NewWithConfig(handlerRequest, configRequest)
    log.ResponseHandler = sloggin.NewWithConfig(handlerResponse, configResponse)
}

func (log *Logger) Close() {
    defer log.Request.Close()
    defer log.Response.Close()
}
