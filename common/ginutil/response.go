package ginutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg,omitempty"`
	Code int         `json:"code"`
}

// Response body
func Response(ctx *gin.Context, err error, data interface{}) {
	if data == nil {
		if err != nil {
			data = "failed"
		} else {
			data = "succeed"
		}
	}
	resp := &responseData{
		Data: data,
	}
	if err != nil {
		resp.Code = 1
		resp.Msg = err.Error()
	}
	ctx.JSON(http.StatusOK, resp)
}
