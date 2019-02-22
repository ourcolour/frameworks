package controllers

import (
	"github.com/gin-gonic/gin"
	. "github.com/ourcolour/frameworks/entities/responses"
	"net/http"
)

func JsonWithStatusCode(c *gin.Context, data interface{}, err error, httpStatusCode int) {
	resp := MakeJsonResponse(data, err)

	c.JSON(httpStatusCode, resp)
}

func Json(c *gin.Context, data interface{}, err error) {
	resp := MakeJsonResponse(data, err)

	var httpStatusCode int
	switch resp.Code {
	case MSG_CODE_WARNING:
		httpStatusCode = http.StatusAccepted
		break
	case MSG_CODE_ERROR:
		httpStatusCode = http.StatusBadRequest
		break
	default:
		httpStatusCode = http.StatusOK
	}

	c.JSON(httpStatusCode, resp)
}
