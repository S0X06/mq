package handler

import (
	// "fmt"
	"net/http"
	// "user/router/middleware"

	"mq/utils"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(code int, c *gin.Context, data interface{}) {

	code, message := utils.GetCode(code)

	if msg, ok := data.(string); ok {
		message = msg
	}

	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	c.JSON(http.StatusOK, response)

}

func Index(c *gin.Context) {

}
