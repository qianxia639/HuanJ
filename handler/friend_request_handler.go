package handler

import (
	"HuanJ/config"
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sendFriendReq struct {
	ToUserId    int32  `json:"to_user_id" binding:"required"`
	RequestDesc string `json:"request_desc"`
}

// 发送好友申请
func (handler *Handler) sendFriendRequest(ctx *gin.Context) {
	var req sendFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ParamsError(ctx)
		return
	}

	// 不能添加自己
	if req.ToUserId == handler.CurrentUserInfo.ID {
		handler.ParamsError(ctx, "不能添加自己")
		return
	}

	// 检查用户是否存在
	u, _ := handler.Store.GetUserById(ctx, req.ToUserId)
	if u.ID < 1 {
		handler.ParamsError(ctx, "用户不存在")
		return
	}

	// 检查是否已是好友
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
		FromUserID: handler.CurrentUserInfo.ID,
		ToUserID:   req.ToUserId,
	}); exists {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	// 检查是否已有申请
	if fr, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{}); fr.ID > 0 {
		handler.Success(ctx, "已申请, 无需重复申请")
		return
	}

	// 添加申请记录
	if err := handler.Store.CreateFriendRequest(ctx, &db.CreateFriendRequestParams{
		FromUserID:  handler.CurrentUserInfo.ID,
		ToUserID:    req.ToUserId,
		RequestDesc: req.RequestDesc,
	}); err != nil {
		logs.Errorf("添加申请记录失败: %v", err)
		handler.ServerError(ctx)
		return
	}

	handler.Success(ctx, "申请成功")
}

// 获取待处理的好友申请
func (handler *Handler) listFriendRequest(ctx *gin.Context) {
	// userId := handler.CurrentUserInfo.ID

	// SELECT * FROM friend_requests WHERE to_user_id = {userId}
	// handler.Store.GetFriendRequest()
}

type ProcessFriendRequest struct {
	FromUserId int32 `json:"from_user_id" binding:"required"`
	// ToUserId   int32  `json:"to_user_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=accept reject"`
	Note   string `json:"note"`
}

// 处理好友申请
func (handler *Handler) processFriendRequest(ctx *gin.Context) {
	var req ProcessFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ParamsError(ctx)
		return
	}

	fr, err := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: req.FromUserId,
		ToUserID:   handler.CurrentUserInfo.ID,
	})
	if fr.ID < 0 {
		logs.Errorf("get friend request error: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不存在的申请记录"})
		return
	}

	// if time.Now().After(fr.RequestedAt) {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请求已过期"})
	// 	return
	// }

	switch req.Action {
	case "accept":
		err := handler.acceptedUserProcess(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	case "reject":
		err = handler.Store.UpdateFriendRequest(ctx, &db.UpdateFriendRequestParams{
			FromUserID: req.FromUserId,
			ToUserID:   handler.CurrentUserInfo.ID,
			Status:     config.Rejected,
		})
		if err != nil {
			handler.ServerError(ctx)
			return
		}
	default:
		handler.ParamsError(ctx)
		return
	}

}

// 同意申请
func (handler *Handler) acceptedUserProcess(ctx context.Context, req ProcessFriendRequest) error {
	args := db.FriendRequestTxParams{
		FromUserId: req.FromUserId,
		ToUserId:   handler.CurrentUserInfo.ID,
		Status:     config.Accepted,
		FromNote:   req.Note,
		ToNote:     handler.CurrentUserInfo.Nickname,
	}
	return handler.Store.FriendRequestTx(ctx, args)
}
