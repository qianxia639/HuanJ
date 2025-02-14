package handler

import (
	db "Ice/db/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createGroupRequest struct {
	GroupName   string `json:"group_name" binding:"required"`
	Description string `json:"description"`
}

func (handler *Handler) createGroup(ctx *gin.Context) {
	var req createGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 判断群组是否已存在
	group, _ := handler.Queries.GetGroup(ctx, req.GroupName)
	if group.ID > 0 {
		Error(ctx, http.StatusBadRequest, "群组名存在")
		return
	}

	// 创建群组
	//  TODO 群组创建成功后应在群组成员表中插入记录
	if err := handler.Queries.CreateGroup(ctx, &db.CreateGroupParams{
		CreatorId:   handler.CurrentUserInfo.ID,
		GroupName:   req.GroupName,
		Description: req.Description,
		UserId:      handler.CurrentUserInfo.ID,
		Role:        Owner,
	}); err != nil {
		Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, nil)
}
