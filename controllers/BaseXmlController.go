package controllers

import (
	"github.com/gin-gonic/gin"
)

func XmlWithStatusCode(c *gin.Context, data interface{}, err error, httpStatusCode int) {
	//resp := MakeJsonResponse(data, err)
	c.Header("Content-Type", "text/xml")
	c.XML(httpStatusCode, data)
}

//
//func Xml(c *gin.Context, data interface{}, err error) {
//	resp := MakeJsonResponse(data, err)
//
//	var httpStatusCode int
//	switch resp.Code {
//	case MSG_CODE_WARNING:
//		httpStatusCode = http.StatusAccepted
//		break
//	case MSG_CODE_ERROR:
//		httpStatusCode = http.StatusBadRequest
//		break
//	default:
//		httpStatusCode = http.StatusOK
//	}
//
//	c.Header("Content-Type", "text/xml")
//	c.XML(httpStatusCode, data)
//}
