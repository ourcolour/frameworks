package controllers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func MustGetRequestBody(c *gin.Context) []byte {
	result, err := GetRequestBody(c)

	if nil != err {
		Json(c, result, err)
		return nil
	}

	return result
}

func GetRequestBody(c *gin.Context) ([]byte, error) {
	var (
		result []byte
		err    error
	)

	// 参数
	result, err = ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	return result, err
}
