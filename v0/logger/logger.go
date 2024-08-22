package logger

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Logger struct {
	Request  *os.File
	Response *os.File

	RequestHandler  gin.HandlerFunc
	ResponseHandler gin.HandlerFunc
}

var LoggerInst Logger

func (logger *Logger) Init() {

	var err error
	logger.Request, err = os.Create(os.Getenv("REQUEST_LOG") + ".request.log")
	if err != nil {
		log.Println("Error creating log file for requests", err)
	}

	logger.Response, err = os.Create(os.Getenv("RESPONSE_LOG") + ".response.log")
	if err != nil {
		log.Println("Error creating log file for response", err)
	}

	configRequest := sloggin.Config{
		WithRequestBody:    true,
		WithResponseBody:   false,
		WithRequestHeader:  true,
		WithResponseHeader: false,
	}

	configResponse := sloggin.Config{
		WithRequestBody:    false,
		WithResponseBody:   true,
		WithRequestHeader:  false,
		WithResponseHeader: true,
	}

	handlerRequest := slog.New(slog.NewJSONHandler(io.MultiWriter(logger.Request), nil))
	handlerResponse := slog.New(slog.NewJSONHandler(io.MultiWriter(logger.Response), nil))

	logger.RequestHandler = sloggin.NewWithConfig(handlerRequest, configRequest)
	logger.ResponseHandler = sloggin.NewWithConfig(handlerResponse, configResponse)
}

func (logger *Logger) Close() {
	defer logger.Request.Close()
	defer logger.Response.Close()
}
