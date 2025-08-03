package handler

import (
	"HuanJ/config"
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
		UserID:   currentUserId,
		FriendID: req.ToUserId,
	}); exists {
		ctx.JSON(http.StatusOK, "已经是好友")
		return
	}

	// 检查申请是否存在(当前用户发送给目标用户的)
	exists, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: currentUserId,
		ToUserID:   req.ToUserId,
	})
	if exists {
		handler.Success(ctx, "已发送申请, 等待对方处理")
		return
	}

	// 检查是否有来自目标用户的待处理申请(反向申请)
	incomingReqExists, _ := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: req.ToUserId,
		ToUserID:   currentUserId,
	})

	// 如果存在反向申请, 则自动接受
	if incomingReqExists {
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

func (handler *Handler) sendFriendRequestV2(ctx *gin.Context) {
	var req sendFriendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ParamsError(ctx)
		return
	}

	currentUserInfo := ctx.MustGet("current_user_info").(db.LoginUserInfo)

	// 不能添加自己
	if req.ToUserId == currentUserInfo.ID {
		handler.ParamsError(ctx, "不能添加自己")
		return
	}

	// 使用缓存防止重复提交
	cacheKey := fmt.Sprintf("friend_req:%d:%d", currentUserInfo.ID, req.ToUserId)
	if exists, _ := handler.RedisClient.SetNX(ctx, cacheKey, 1, 5*time.Second).Result(); exists {
		handler.Success(ctx, "请求处理中，请勿重复提交")
		return
	}

	// 并行查询目标用户和好友关系
	var wg sync.WaitGroup
	wg.Add(3)
	g, gCtx := errgroup.WithContext(ctx)
	// 检查目标用户是否存在
	var targetUser db.User
	g.Go(func() error {
		user, _ := handler.Store.GetUserById(gCtx, req.ToUserId)
		if user.ID < 1 {
			return fmt.Errorf("用户不存在")
		}
		targetUser = user
		return nil
	})

	// 检查是否已是好友
	var isFriend bool
	g.Go(func() error {
		isFriend, _ = handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
			UserID:   currentUserInfo.ID,
			FriendID: req.ToUserId,
		})
		return nil
	})

	// 单次查询获取双向申请记录
	var forwardReq *db.FriendRequest
	var reverseReq *db.FriendRequest
	g.Go(func() error {
		minId, maxId := currentUserInfo.ID, req.ToUserId
		if minId > maxId {
			minId, maxId = req.ToUserId, currentUserInfo.ID
		}
		// requests, err := handler.Store.GetMutualFriendRequests(ctx, &db.GetMutualFriendRequestsParams{
		// 	UserA: currentUserId,
		// 	UserB: req.ToUserId,
		// })
		// if err == nil && len(requests) > 0 {
		// 	for _, r := range requests {
		// 		if r.FromUserID == currentUserId {
		// 			forwardReq = &r
		// 		} else {
		// 			reverseReq = &r
		// 		}
		// 	}
		// }
		// 分离双向申请记录
		// for i := range requests {
		// 	switch requests[i].FromUserID {
		// 	case currentUserInfo.ID:
		// 		forwardReq = &requests[i]
		// 	case req.ToUserId:
		// 		reverseReq = &requests[i]
		// 	}
		// }
		return nil
	})

	if err := g.Wait(); err != nil {
		zap.L().Error("查询失败", zap.Error(err))
		handler.ServerError(ctx)
		return
	}

	// 校验用户是否存在
	if targetUser.ID < 1 {
		handler.ParamsError(ctx, "用户不存在")
		return
	}

	// 检查是否已是好友
	if isFriend {
		handler.Success(ctx, "已经是好友")
		return
	}
	// 检查申请是否存在(当前用户发送给目标用户的)
	if forwardReq != nil && forwardReq.Status == config.Pending {
		handler.Success(ctx, "已发送申请, 等待对方处理")
		return
	}

	// 处理反向申请(自动接受)
	if reverseReq != nil && reverseReq.Status == config.Pending {
		if err := handler.Store.FriendRequestTx(ctx, &db.FriendRequestTxParams{
			Status:     config.Accepted,    // 同意申请
			FromUserId: req.ToUserId,       // 申请发起者是目标用户
			ToUserId:   currentUserInfo.ID, // 当前用户是接收者
			FromNote:   currentUserInfo.Nickname,
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
		FromUserID:  currentUserInfo.ID,
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

	exists, err := handler.Store.GetFriendRequest(ctx, &db.GetFriendRequestParams{
		FromUserID: req.FromUserId,
		ToUserID:   handler.CurrentUserInfo.ID,
	})
	if !exists {
		logs.Errorf("get friend request error: %v", err.Error())
		handler.Error(ctx, http.StatusNotFound, "不存在的申请记录")
		return
	}

	// 校验当前用户是否是接收者
	// if fr.ToUserID != handler.CurrentUserInfo.ID {
	// 	logs.Errorf("无权处理此请求, senderId: %d, receiveId: %d", fr.ToUserID, handler.CurrentUserInfo.ID)
	// 	handler.ParamsError(ctx)
	// 	return
	// }

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
