package repository

import (
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"time"

	"gorm.io/gorm"
)

// 获取讨论列表
func SelectDiscussList(db *gorm.DB, page int, itemsPerPage int) ([]structs.DiscussRequest, int) {
	var discussRequests []structs.DiscussRequest
	var total int64

	db.Table("Discussions").
		Select("Discussions.Did, Discussions.Title, Users.Username, Discussions.Create_at").
		Joins("JOIN Users ON Discussions.Uid = Users.Uid").
		Order("Discussions.Create_at desc").Count(&total).
		Offset((page - 1) * itemsPerPage).Limit(itemsPerPage).Find(&discussRequests)

	return discussRequests, int(total)
}

// 获取指定Did讨论及User表中发帖人的头像
func SelectDiscussionByDid(db *gorm.DB, Did int) structs.DiscsInfoRequest {
	var discussion structs.DiscsInfoRequest
	db.Table("Discussions").
		Select("Discussions.Did, Discussions.Title, Discussions.Content, Discussions.Create_at, Users.Uid,Users.Username, Users.Avatar").
		Joins("JOIN Users ON Discussions.Uid = Users.Uid").
		Where("Discussions.Did = ?", Did).First(&discussion)
	return discussion
}

// 添加讨论
func AddDiscussion(db *gorm.DB, title, content string, uid int) bool {
	if title == "" || content == "" {
		return false
	}
	err := db.Table("Discussions").
		Create(&tables.Discussion{UserId: uid, Title: title, Content: content, CreatedAt: time.Now()}).Error
	return err == nil
}

// 删除讨论
func DelDiscussion(db *gorm.DB, Did int) bool {
	err := db.Table("Discussions").Where("Did = ?", Did).Delete(&tables.Discussion{}).Error
	return err == nil
}

// 添加评论
func AddComment(db *gorm.DB, content string, did, uid int, profanity bool) bool {
	return db.Table("Comments").Create(&tables.Comment{Id: did, UserId: uid, Content: content, CreatedAt: time.Now(), Profanity: profanity}).Error == nil
}

// 获取指定讨论ID的所有评论信息
func SelectCommentsByDid(db *gorm.DB, Did int) []structs.CommentRequest {
	var comments []structs.CommentRequest
	db.Table("Comments").
		Select("Comments.Cid, Comments.Did, Comments.Content, Comments.Create_at, Users.Uid,Users.Username, Users.Avatar,Comments.Profanity").
		Joins("JOIN Users ON Comments.Uid = Users.Uid").
		Order("create_at desc").
		Where("Comments.Did = ?", Did).Find(&comments)
	return comments
}

// 删除指定评论
func DeleteCommentByCid(db *gorm.DB, Cid int) bool {
	return db.Table("Comments").Where("Cid = ?", Cid).Delete(&tables.Comment{}).Error == nil
}
