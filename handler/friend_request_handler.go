package handler

import (
	db "Rejuv/db/sqlc"
	"Rejuv/logs"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createFriendReq struct {
	ToUserId    int32  `json:"to_user_id" binding:"required"`
	RequestDesc string `json:"request_desc"`
}

// 添加好友申请
func (handler *Handler) createFriendRequest(ctx *gin.Context) {
	var req createFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 校验是否是自己申请
	if req.ToUserId == handler.CurrentUserInfo.ID {
		logs.Errorf("userId: %d, friendId: %d\n", handler.CurrentUserInfo.ID, req.ToUserId)
		Error(ctx, http.StatusUnauthorized, "不能添加自己")
		return
	}

	// 检查请求者是否存在
	u, _ := handler.Store.GetUserById(ctx, req.ToUserId)
	if u.ID == 0 {
		Error(ctx, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 检查是否存在待处理的请求
	if fr, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{}); fr.ID > 0 {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	// 检查是否已经是好友
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
		SenderID:   handler.CurrentUserInfo.ID,
		ReceiverID: req.ToUserId,
	}); exists {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	if err := handler.Store.CreateFriendRequest(ctx, &db.CreateFriendRequestParams{
		SenderID:    handler.CurrentUserInfo.ID,
		ReceiverID:  req.ToUserId,
		RequestDesc: req.RequestDesc,
	}); err != nil {
		logs.Error(err)
		Error(ctx, http.StatusInternalServerError, "申请失败")
		return
	}

	Success(ctx, "申请成功")
}

func (handler *Handler) listFriendRequest(ctx *gin.Context) {
	// userId := handler.CurrentUserInfo.ID

	// SELECT * FROM friend_requests WHERE to_user_id = {userId}
	// handler.Store.GetFriendRequest()
}

type ProcessFriendRequest struct {
	FromUserId int32  `json:"from_user_id" binding:"required"`
	Action     string `json:"status" binding:"required,oneof=accept reject"`
	Note       string `json:"note"`
}

func (handler *Handler) processFriendRequest(ctx *gin.Context) {
	var req ProcessFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	fr, err := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		SenderID:   req.FromUserId,
		ReceiverID: handler.CurrentUserInfo.ID,
	})
	if fr.ID < 0 {
		logs.Errorf("get friend request error: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不存在的申请记录"})
		return
	}

	if time.Now().After(fr.ExpiredAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请求已过期"})
		return
	}

	switch req.Action {
	case "accept":
		err := handler.acceptedUserProcess(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	case "reject":
		err = handler.rejectedUserProcess(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

}

// 同意申请
func (handler *Handler) acceptedUserProcess(ctx context.Context, req ProcessFriendRequest) error {
	args := db.FriendRequestTxParams{
		FromUserId: req.FromUserId,
		ToUserId:   handler.CurrentUserInfo.ID,
		Status:     Accepted,
		FromNote:   req.Note,
		ToNote:     handler.CurrentUserInfo.Nickname,
	}
	return handler.Store.FriendRequestTx(ctx, args)
}

// 拒绝申请
func (handler *Handler) rejectedUserProcess(ctx context.Context, req ProcessFriendRequest) error {
	return handler.Store.UpdateFriendRequest(ctx, &db.UpdateFriendRequestParams{
		SenderID:   req.FromUserId,
		ReceiverID: handler.CurrentUserInfo.ID,
		Status:     Rejected,
	})
}
