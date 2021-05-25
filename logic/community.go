/**
 * @Author: Robby
 * @File name: community.go
 * @Create date: 2021-05-24
 * @Function:
 **/

package logic

import (
	mysqlconnect "jiaoshoujia/dao/mysql"
	"jiaoshoujia/models"
)

// GetCommunityList ：这里拿到列表结果，直接返回
func GetCommunityList() (data []*models.Community, err error) {
	data, err = mysqlconnect.GetCommunityList()
	return
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {

	data, err = mysqlconnect.GetCommunityDetailById(id)
	return

}
