/**
 * @Author: Robby
 * @File name: community.go
 * @Create date: 2021-05-24
 * @Function:
 **/

package models

import "time"

type Community struct {
	// json的反序列化字段名使用id，数据库字段名使用community_id
	Id   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	Id           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
