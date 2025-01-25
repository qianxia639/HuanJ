package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createdFriendRequest struct {
	FromUserId  uint32 `json:"from_user_id" binding:"required"` // 申请者Id
	ToUserId    uint32 `json:"to_user_id" binding:"required"`   // 接收者Id
	Description string `json:"description" binding:"required"`  // 申请描述
}

func (h *Handler) createdFriend(ctx *gin.Context) {

	// 判断是否存在申请记录
	// 1.存在A申请B的记录直接返回
	// 2.存在A申请B的记录，B又想A申请则直接同意
	var req createdFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.FromUserId == req.ToUserId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "不能添加自己为好友"})
		return
	}

	if u, _ := h.Queries.GetUserById(ctx, req.ToUserId); u.ID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户不存在"})
		return
	}

	// 身份校验
	// auth := ctx.Request.Header.Get(authorizationHeader)
	// fields := strings.Fields(auth)
	// payload, err := h.Token.VerifyToken(fields[1])
	// if err != nil {
	// 	Error(ctx, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// if payload.Username != h.CurrentUserInfo.Username {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"messagge": "权限不足"})
	// 	return
	// }
	if req.FromUserId != h.CurrentUserInfo.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"messagge": "权限不足",
			"from_user_id": req.FromUserId, "id": h.CurrentUserInfo.ID})
		return
	}

	// 判断是否已经是好友
	if i := h.Queries.ExistsFriend(ctx, req.FromUserId, req.ToUserId, ACCEPTED); i > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "已经是好友"})
		return
	}

	// 判断关系是否存在
	// 如果A申请B存在，则直接返回
	// 如果A申请B存爱且B又申请A，则B同意A的申请
	if i := h.Queries.ExistsFriend(ctx, req.FromUserId, req.ToUserId, PENDING); i > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "关系存在"})
		return
	}

	// 判断是否有来自对方的申请
	// 存在则同意
	if i := h.Queries.ExistsFriend(ctx, req.ToUserId, req.FromUserId, PENDING); i > 0 {
		if err := h.Queries.AddFriendTx(ctx, req.FromUserId, req.ToUserId); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
		return
	}

	if err := h.Queries.AddFriendRequest(ctx, req.FromUserId, req.ToUserId, ""); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
}

func (h *Handler) getFriends(ctx *gin.Context) {

	query := ctx.Query("user_id")

	userId, err := strconv.Atoi(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	friends, _ := h.Queries.GetFriendAll(ctx, uint32(userId))

	ctx.JSON(http.StatusOK, friends)

}

type deleteFriendRequest struct {
	FromUserId uint32 `json:"from_user_id" binding:"required"`
	ToUserId   uint32 `json:"to_user_id" binding:"required"`
}

func (h *Handler) deleteFriend(ctx *gin.Context) {
	var req deleteFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 无法删除自己
	if req.FromUserId == req.ToUserId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无法删除自己"})
		return
	}

	if req.FromUserId != h.CurrentUserInfo.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"messagge": "权限不足",
			"from_user_id": req.FromUserId, "id": h.CurrentUserInfo.ID})
		return
	}

	// 判断要删除用户是否存在
	if u, err := h.Queries.GetUserById(ctx, req.ToUserId); u.ID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "用户不存在"})
		return
	}

	// 判断是否是好友
	if i := h.Queries.ExistsFriend(ctx, req.FromUserId, req.ToUserId, ACCEPTED); i < 1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "非好友无法删除"})
		return
	}

	// 删除
	err := h.Queries.DeleteFriend(ctx, req.FromUserId, req.ToUserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
}
