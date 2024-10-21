package handler

import (
	"Dandelion/token"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createdFriendRequest struct {
	UserId   uint32 `json:"user_id" binding:"required"`
	FriendId uint32 `json:"friend_id" binding:"required"`
}

func (h *Handler) createdFriend(ctx *gin.Context) {
	var req createdFriendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserId == req.FriendId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "不能添加自己为好友"})
		return
	}

	// 身份校验
	data, exist := ctx.Get(authorizationPayloadKey)
	if !exist {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not user"})
		return
	}
	payload := data.(*token.Payload)

	user, _ := h.Queries.GetUser(ctx, payload.Username)
	fmt.Printf("user.ID: %v\n", user.ID)
	if req.UserId != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "身份认证失败"})
		return
	}

	// 判断关系是否存在
	friend, _ := h.Queries.GetFriend(ctx, req.UserId, req.FriendId)
	if friend.ID > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "关系存在"})
		return
	}

	err := h.Queries.AddFriend(ctx, req.UserId, req.FriendId)
	if err != nil {
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
