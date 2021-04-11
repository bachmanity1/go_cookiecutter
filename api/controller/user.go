package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func newHTTPUserHandler(eg *echo.Group, handler *HTTPHandler) {
	// Prefix : /api/v1/user
	ueg := eg.Group("/user")
	ueg.GET("/:uid", handler.GetUserByID)
}

// GetUserByID ...
func (h *HTTPHandler) GetUserByID(c echo.Context) (err error) {
	ctx := c.Request().Context()

	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		mlog.Errorw("GetUserByID", "error", err)
		return response(c, http.StatusBadRequest, "Invalid Path Param")
	}

	user, err := h.uService.GetUserByID(ctx, uid)
	if err != nil {
		mlog.Errorw("GetUserByID", "error", err)
		return response(c, http.StatusInternalServerError, err.Error())
	}

	return response(c, http.StatusOK, "GetUserByID OK", user)
}
