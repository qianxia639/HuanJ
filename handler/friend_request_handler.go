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

	currentUserId := handler.CurrentUserInfo.ID

	// 不能添加自己
	if req.ToUserId == currentUserId {
		handler.ParamsError(ctx, "不能添加自己")
		return
	}

	// 检查目标用户是否存在
	targetUser, _ := handler.Store.GetUserById(ctx, req.ToUserId)
	if targetUser.ID < 1 {
		handler.ParamsError(ctx, "用户不存在")
		return
	}

	// 检查是否已是好友
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
		FromUserID: currentUserId,
		ToUserID:   req.ToUserId,
	}); exists {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	// 检查是否已有申请(当前用户发送给目标用户的)
	outgoingReq, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: currentUserId,
		ToUserID:   req.ToUserId,
	})
	if outgoingReq.ID > 0 {
		handler.Success(ctx, "已发送申请, 等待对方处理")
		return
	}

	// 检查是否有来自目标用户的待处理申请(反向申请)
	incomingReq, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: req.ToUserId,
		ToUserID:   currentUserId,
	})

	// 如果存在反向申请, 则自动接受
	if incomingReq.ID > 0 {
		if err := handler.Store.FriendRequestTx(ctx, &db.FriendRequestTxParams{
			Status:     config.Accepted, // 同意申请
			FromUserId: req.ToUserId,    // 申请发起者是目标用户
			ToUserId:   currentUserId,   // 当前用户是接收者
			FromNote:   handler.CurrentUserInfo.Nickname,
			ToNote:     targetUser.Nickname,
		}); err != nil {
			logs.Errorf("自动接受好友申请失败: %v", err)
			handler.ServerError(ctx)
			return
		}
		handler.Success(ctx, "success")
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

type processFriendRequest struct {
	FromUserId int32 `json:"from_user_id" binding:"required"`
	// ToUserId   int32  `json:"to_user_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=accept reject"`
	Remark string `json:"remark"`
}

// 处理好友申请
func (handler *Handler) processFriendRequest(ctx *gin.Context) {
	var req processFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ParamsError(ctx)
		return
	}

	fr, err := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: req.FromUserId,
		ToUserID:   handler.CurrentUserInfo.ID,
	})
	if fr.ID <= 0 {
		logs.Errorf("get friend request error: %v", err.Error())
		handler.Error(ctx, http.StatusNotFound, "不存在的申请记录")
		return
	}

	// 校验当前用户是否是接收者
	if fr.ToUserID != handler.CurrentUserInfo.ID {
		logs.Errorf("无权处理此请求, senderId: %d, receiveId: %d", fr.ToUserID, handler.CurrentUserInfo.ID)
		handler.ParamsError(ctx)
		return
	}

	switch req.Action {
	case "accept":
		err := handler.acceptedUserProcess(ctx, req)
		if err != nil {
			logs.Errorf("accepted user error: %v\n", err)
			handler.ServerError(ctx)
			return
		}
	case "reject":
		err = handler.Store.UpdateFriendRequest(ctx, &db.UpdateFriendRequestParams{
			FromUserID: req.FromUserId,
			ToUserID:   handler.CurrentUserInfo.ID,
			Status:     config.Rejected,
		})
		if err != nil {
			logs.Errorf("rejected user error: %v\n", err)
			handler.ServerError(ctx)
			return
		}
	default:
		handler.ParamsError(ctx)
		return
	}

}

// 同意申请
func (handler *Handler) acceptedUserProcess(ctx context.Context, req processFriendRequest) error {
	args := &db.FriendRequestTxParams{
		FromUserId: req.FromUserId,
		ToUserId:   handler.CurrentUserInfo.ID,
		Status:     config.Accepted,
		FromNote:   req.Remark,
		ToNote:     handler.CurrentUserInfo.Nickname,
	}
	return handler.Store.FriendRequestTx(ctx, args)
}
