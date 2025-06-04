package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusinessError struct {
	Code    int
	Message string
	Err     error
}

func (e *BusinessError) Error() string {
	return e.Message
}

type JSONErrorResponse struct {
	Error   bool
	Message string
	Code    int
}

type JSONResponse struct {
	Error   bool
	Message string
	Data    interface{}
}

func ResponseJson(c *gin.Context, err bool, msg string, data interface{}, code ...int) {
	if len(code) > 0 {
		c.JSON(http.StatusOK, JSONErrorResponse{
			Error:   err,
			Message: msg,
			Code:    code[0],
		})
		return
	}
	c.JSON(http.StatusOK, JSONResponse{
		Error:   err,
		Message: msg,
		Data:    data,
	})
}
