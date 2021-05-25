/**
 * @Author: Robby
 * @File name: community.go
 * @Create date: 2021-05-24
 * @Function:
 **/

package mysqlconnect

import (
	"database/sql"
	"fmt"
	"jiaoshoujia/models"

	"go.uber.org/zap"
)

// GetCommunityList 查询记录列表，映射到communityList中
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	// 这里由于返回申明的是communityList，不是地址，所以这里必须填写&communityList
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		// 如果没有查询到任何数据
		if err == sql.ErrNoRows {
			zap.L().Warn("没有查询到数据", zap.Any("sql", "select community_id, community_name from community"))
			err = nil
		}
	}
	return
}

// GetCommunityDetailById  基于ID，获取详情数据
func GetCommunityDetailById(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	err = db.Get(communityDetail, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("没有查询到数据", zap.Any("sql", fmt.Sprintf("select community_id, community_name, introduction, create_time from community where community_id = ?", id)))
			err = ErrorInvalidId
		}
	}
	return
}
