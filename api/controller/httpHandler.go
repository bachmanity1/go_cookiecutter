package controller

import (
	"pandita/conf"
	"pandita/service"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	pandita  *conf.ViperConfig
	uService service.UserService
}

func newHTTPHandler(eg *echo.Group,
	pandita *conf.ViperConfig,
	uService service.UserService) {
	handler := &HTTPHandler{pandita, uService}

	userGroup := eg.Group("/user")
	newHTTPUserHandler(userGroup, handler)
}
