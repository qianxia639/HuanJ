package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	OK   bool        `json:"ok"`             // 请求是否成功
	Code int         `json:"code.omitempty"` // 业务状态码(自定义)
	Msg  string      `json:"msg,omitempty"`  // 错误信息
	Data interface{} `json:"data,omitempty"` // 数据信息
}

type ResiltPage struct {
	PageNo   int         `json:"page_on,omitempty"`   // 页号码
	PageSize int         `json:"page_size,omitempty"` // 页大小
	Total    int         `json:"total,omitempty"`     // 数据总条目
	List     interface{} `json:"list"`                // 数据信息
}

func Obj(ctx *gin.Context, data interface{}) {
	result := Result{OK: true, Data: data}
	ctx.JSON(http.StatusOK, result)

}

func Success(ctx *gin.Context, msg string) {
	result := Result{OK: true, Msg: msg}
	ctx.JSON(http.StatusOK, result)
}

func Error(ctx *gin.Context, httpCode int, msg string) {
	result := Result{OK: false, Msg: msg}
	ctx.JSON(httpCode, result)
}

func ServerError(ctx *gin.Context) {
	result := Result{OK: false, Msg: "服务异常"}
	ctx.JSON(http.StatusInternalServerError, result)
	return
}

func ParamsError(ctx *gin.Context, msg ...string) {
	var result Result
	if len(msg) < 1 {
		result = Result{OK: false, Msg: "参数错误"}
	} else {
		result = Result{OK: false, Msg: msg[0]}
	}

	// result := Result{OK: false, Msg: "参数错误"}
	ctx.JSON(http.StatusBadRequest, result)
}
