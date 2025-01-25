package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createFriendRecordRequest struct {
	UserId   uint32 `json:"user_id" binding:"required"`
	FriendId uint32 `json:"friend_id" binding:"required"`
}

func createFriendRecord(ctx *gin.Context) {
	var req createFriendRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
}
