package handler

import (
	db "Ice/db/sqlc"
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
	group, _ := handler.Store.GetGroup(ctx, req.GroupName)
	if group.ID > 0 {
		Error(ctx, http.StatusBadRequest, "群组名存在")
		return
	}

	// 创建群组
	// 创建成功后并将创建者信息写入群员表
	// result, err := handler.Store.CreateGroupTx(ctx, db.CreateGroupTxParams{
	// 	CreateGroupParams: db.CreateGroupParams{
	// 		GroupName:   req.GroupName,
	// 		CreatorID:   handler.CurrentUserInfo.ID,
	// 		Description: req.Description,
	// 	},
	// 	AfterCreate: func(group db.Group) error {
	// 		_, err := handler.Store.CreateGroupMember(ctx, &db.CreateGroupMemberParams{
	// 			GroupID: group.ID,
	// 			UserID:  handler.CurrentUserInfo.ID,
	// 			Role:    GroupOwner,
	// 			Agreed:  true,
	// 		})
	// 		return err
	// 	},
	// })

	result, err := handler.Store.CreateGroupTx(ctx, db.CreateGroupTxParams{
		CreateGroupParams: db.CreateGroupParams{
			GroupName:   req.GroupName,
			CreatorID:   handler.CurrentUserInfo.ID,
			Description: req.Description,
		},
		UserId: handler.CurrentUserInfo.ID,
		Role:   GroupOwner,
		Agreed: true,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	Success(ctx, result)
}
