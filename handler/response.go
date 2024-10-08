package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	OK      bool        `json:"ok"`                // 请求是否成功
	Message string      `json:"message,omitempty"` // 错误信息
	Data    interface{} `json:"data,omitempty"`    // 数据信息
}

type ResultList struct {
	List     interface{} `json:"list"`                // 数组数据
	Total    *int64      `json:"total,omitempty"`     // 数据区总条数
	PageNo   *int        `json:"page_no,omitempty"`   // 页号码
	PageSize *int        `json:"page_size,omitempty"` // 页大小
}

func Success(ctx *gin.Context, data interface{}) {
	result := Result{OK: true, Message: "Successfully", Data: data}
	ctx.JSON(http.StatusOK, result)
}

func Error(ctx *gin.Context, code int, message string) {
	result := Result{OK: false, Message: message}
	ctx.JSON(code, result)
}
