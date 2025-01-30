package handler

import (
	"Ice/internal/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createFriendRecordRequest struct {
	// FromUserId  uint32 `json:"from_user_id" binding:"required"`
	ToUserId    uint32 `json:"to_user_id" binding:"required"`
	RequestDesc string `json:"request_desc" binding:"required"`
}

func (handler *Handler) createFriendRequest(ctx *gin.Context) {
	var req createFriendRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if req.ToUserId == handler.CurrentUserInfo.ID {
		logs.Errorf("userId: %d, friendId: %d\n", handler.CurrentUserInfo.ID, req.ToUserId)
		Error(ctx, http.StatusUnauthorized, "不能添加自己")
		return
	}

	if handler.Queries.ExistsFriendRecord(ctx, handler.CurrentUserInfo.ID, req.ToUserId) > 0 {
		Error(ctx, http.StatusInternalServerError, "已发送申请")
		return
	}

	if err := handler.Queries.AddFriendRequest(ctx, handler.CurrentUserInfo.ID, req.ToUserId, req.RequestDesc); err != nil {
		logs.Error(err)
		Error(ctx, http.StatusInternalServerError, "申请失败")
		return
	}

	Success(ctx, "申请成功")
}

type acceptFriendRequest struct {
	Id uint32 `json:"id" binding:"required"`
}

func (handler *Handler) acceptFriendRequest(ctx *gin.Context) {
	var req acceptFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 判断要处理的记录是否存在
	if count := handler.Queries.ExistsFriendRequest(ctx, req.Id); count < 1 {
		Error(ctx, http.StatusInternalServerError, "记录不存在")
		return
	}

	fr, err := handler.Queries.GetFriendRequest(ctx, req.Id, 1)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = handler.Queries.InsertAcceptFriendRequestTx(ctx, req.Id, fr.FromUserId, fr.ToUserId)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, "Successfully")
}

func (handler *Handler) rejectFriendRequest(ctx *gin.Context) {
	var req acceptFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 判断要处理的记录是否存在
	if count := handler.Queries.ExistsFriendRequest(ctx, req.Id); count < 1 {
		Error(ctx, http.StatusInternalServerError, "记录不存在")
		return
	}

	err := handler.Queries.UpdateFriendRequest(ctx, req.Id)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, "Successfully")
}
