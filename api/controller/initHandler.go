package controller

import (
	"fmt"
	"net/http"
	mw "pandita/api/middleware"
	"pandita/conf"
	repo "pandita/repository"
	"pandita/service"
	"pandita/util"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/swag/example/basic/docs"
	"gorm.io/gorm"
)

type (
	// PanditaStatus for common response status
	PanditaStatus struct {
		TRID       string      `json:"trID" example:"20200213052007345858"`
		ResultCode string      `json:"resultCode" example:"0000"`
		ResultMsg  string      `json:"resultMsg" example:"Request OK"`
		ResultData interface{} `json:"resultData,omitempty"`
	}
)

var mlog *util.MLogger

func init() {
	mlog, _ = util.InitLog("controller", "devel")
}

// InitHandler ...
func InitHandler(pandita *conf.ViperConfig, e *echo.Echo, db *gorm.DB) (err error) {

	mlog, _ = util.InitLog("controller", pandita.GetString("loglevel"))
	// timeout := time.Duration(pandita.GetInt("timeout")) * time.Second

	// Default Group
	//	url := echoSwagger.URL("http://localhost:10811/swagger/swagger.json")
	//e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	docs.SwaggerInfo.Host = pandita.GetString("swagger_host")
	api := e.Group("/api")
	ver := api.Group("/v1")
	ver.Use(mw.TransID())

	// envStr := pandita.GetString("ENV")
	// jwtkey := pandita.GetString("jwt_access_key")
	// if jwtkey == "" {
	// 	mlog.Warnw("InitHandler cannot call Client API because of empty JWTKey")
	// }
	// if envStr != "" && envStr != "prod" {
	// 	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// }
	// e.GET("/*", func(c echo.Context) error { return c.Render(http.StatusOK, "envStr", envStr) })

	timeout := 10 * time.Second
	uRepo := repo.NewGormUserRepository(db)
	uService := service.NewUserService(uRepo, timeout)
	newHTTPHandler(ver, pandita, uService)
	return nil
}

func response(c echo.Context, code int, resMsg string, result ...interface{}) error {
	resCode := "0000"
	if code != http.StatusOK {
		resCode = fmt.Sprintf("1%d", code)
	}

	id, ok := c.Request().Context().Value(util.TransIDKey).(string)
	if !ok {
		id = util.NewID()
	}

	res := PanditaStatus{
		TRID:       id,
		ResultCode: resCode,
		ResultMsg:  resMsg,
	}

	if result != nil {
		res.ResultData = result[0]
	}
	return c.JSON(code, res)
}
