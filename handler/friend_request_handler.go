package handler

import (
	db "Rejuv/db/sqlc"
	"Rejuv/logs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createFriendReq struct {
	FriendId    int32  `json:"friend_id" binding:"required"`
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
	if req.FriendId == handler.CurrentUserInfo.ID {
		logs.Errorf("userId: %d, friendId: %d\n", handler.CurrentUserInfo.ID, req.FriendId)
		Error(ctx, http.StatusUnauthorized, "不能添加自己")
		return
	}

	// 检查请求者是否存在
	u, _ := handler.Store.GetUserById(ctx, req.FriendId)
	if u.ID == 0 {
		Error(ctx, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 检查是否已经是好友
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
		UserID:   handler.CurrentUserInfo.ID,
		FriendID: req.FriendId,
	}); exists {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	// 检查是否存在待处理的请求
	if count, _ := handler.Store.ExistsFriendRequest(ctx, &db.ExistsFriendRequestParams{
		UserID:   handler.CurrentUserInfo.ID,
		FriendID: req.FriendId,
	}); count > 0 {
		Error(ctx, http.StatusInternalServerError, "已存在待处理的请求")
		return
	}

	if err := handler.Store.CreateFriendRequest(ctx, &db.CreateFriendRequestParams{
		UserID:      handler.CurrentUserInfo.ID,
		FriendID:    req.FriendId,
		RequestDesc: req.RequestDesc,
	}); err != nil {
		logs.Error(err)
		Error(ctx, http.StatusInternalServerError, "申请失败")
		return
	}

	Success(ctx, "申请成功")
}

type AcceptedOrRejectFriendRequest struct {
	RequestId int32 `json:"request_id" binding:"required"`
	Status    int8  `json:"status"`
}

func (handler *Handler) pendingProcess(ctx *gin.Context) {

	var req AcceptedOrRejectFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

		return
	}

	if count, _ := handler.Store.ExistsFriendRequest(ctx, &db.ExistsFriendRequestParams{
		UserID:   handler.CurrentUserInfo.ID,
		FriendID: req.RequestId,
	}); count < 1 {
		Error(ctx, http.StatusInternalServerError, "处理请求不存在")
		return
	}

	// err = handler.Queries.AcceptFriendRequest(ctx, int32(requestId), handler.CurrentUserInfo.ID)
	// if err != nil {
	// 	Error(ctx, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// 获取申请记录
	// SELECT * FROM friend_requests WHERE (sender_id = {req.RequestId} AND receiver_id = {currentUserId}) OR (sender_id = {currentUserId} AND receiver_id = {req.RequestId})

	// 检查是否过期
	// if friendRequest.Status != 1 || time.Now().After(friendRequest.ExpiredAt) {"申请已过期"}

	user, _ := handler.Store.GetUserById(ctx, req.RequestId)

	args := db.FriendRequestTxParams{
		UserId:       handler.CurrentUserInfo.ID,
		FriendId:     req.RequestId,
		Status:       Accepted,
		FromNickname: handler.CurrentUserInfo.Nickname,
		ToNickname:   user.Nickname,
	}
	err := handler.Store.FriendRequestTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	switch req.Status {
	case 1: // 同意
	case 2: // 拒绝
	default: // 无效

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
	if count, _ := handler.Store.ExistsFriendRequest(ctx, &db.ExistsFriendRequestParams{
		UserID:   handler.CurrentUserInfo.ID,
		FriendID: int32(requestId),
	}); count < 1 {
		Error(ctx, http.StatusInternalServerError, "处理请求不存在")
		return
	}

	args := &db.UpdateFriendRequestParams{
		UserID:   handler.CurrentUserInfo.ID,
		FriendID: int32(requestId),
		Status:   Rejected,
	}
	err = handler.Store.UpdateFriendRequest(ctx, args)
	if err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, "Successfully")
}
