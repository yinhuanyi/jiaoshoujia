/**
 * @Author: Robby
 * @File name: community.go
 * @Create date: 2021-05-24
 * @Function:
 **/

package controllers

import (
	"jiaoshoujia/logic"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CommunityHandler : 获取社区列表
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList()", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)

}

// CommunityDetailHandler : 社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 这里id必须与路由的:id名称对应
	communityIdStr := c.Param("id")
	communityId, err := strconv.ParseInt(communityIdStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityDetail(communityId)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail(communityId)", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
