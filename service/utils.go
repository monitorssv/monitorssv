package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnOk(c *gin.Context, data interface{}) {
	if data == nil {
		data = successRes
	}

	c.JSON(http.StatusOK, data)
}

func ReturnErr(c *gin.Context, res *Response) {
	c.JSON(http.StatusOK, res)
}
