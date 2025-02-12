package handler

import (
	"Ice/internal/logs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createFriendRecordRequest struct {
	FriendId    int32  `json:"friend_id" binding:"required"`
	RequestDesc string `json:"request_desc"`
}

func (handler *Handler) createFriendRequest(ctx *gin.Context) {
	var req createFriendRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if req.FriendId == handler.CurrentUserInfo.ID {
		logs.Errorf("userId: %d, friendId: %d\n", handler.CurrentUserInfo.ID, req.FriendId)
		Error(ctx, http.StatusUnauthorized, "不能添加自己")
		return
	}

	// 检查请求者是否存在
	u, _ := handler.Queries.GetUserById(ctx, req.FriendId)
	if u.ID == 0 {
		Error(ctx, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 检查是否已经是好友
	if count := handler.Queries.CheckFriendship(ctx, handler.CurrentUserInfo.ID, req.FriendId); count > 0 {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	// 检查是否存在待处理的请求
	if count := handler.Queries.ExistsFriendRequest(ctx, handler.CurrentUserInfo.ID, req.FriendId); count > 0 {
		Error(ctx, http.StatusInternalServerError, "已存在待处理的请求")
		return
	}

	if err := handler.Queries.AddFriendRequest(ctx, handler.CurrentUserInfo.ID, req.FriendId, req.RequestDesc); err != nil {
		logs.Error(err)
		Error(ctx, http.StatusInternalServerError, "申请失败")
		return
	}

	Success(ctx, "申请成功")
}

func (handler *Handler) acceptFriendRequest(ctx *gin.Context) {

	requestId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "Invalid param")
		return
	}

	if count := handler.Queries.ExistsFriendRequest(ctx, handler.CurrentUserInfo.ID, int32(requestId)); count < 1 {
		Error(ctx, http.StatusInternalServerError, "处理请求不存在")
		return
	}

	err = handler.Queries.AcceptFriendRequest(ctx, int32(requestId), handler.CurrentUserInfo.ID)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, "Successfully")
}

func (handler *Handler) rejectFriendRequest(ctx *gin.Context) {
	requestId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		Error(ctx, http.StatusBadRequest, "Invalid param")
		return
	}

	// 判断要处理的记录是否存在
	if count := handler.Queries.ExistsFriendRequest(ctx, handler.CurrentUserInfo.ID, int32(requestId)); count < 1 {
		Error(ctx, http.StatusInternalServerError, "处理请求不存在")
		return
	}

	err = handler.Queries.RejectFriendRequest(ctx, int32(requestId), handler.CurrentUserInfo.ID)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, "Successfully")
}
