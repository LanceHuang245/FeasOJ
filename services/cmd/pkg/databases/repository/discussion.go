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
		Select("Discussions.id, Discussions.title, Users.username, Discussions.created_at").
		Joins("JOIN Users ON Discussions.user_id = Users.id").
		Order("Discussions.created_at desc").Count(&total).
		Offset((page - 1) * itemsPerPage).Limit(itemsPerPage).Find(&discussRequests)

	return discussRequests, int(total)
}

// 获取指定Did讨论及User表中发帖人的头像
func SelectDiscussionByDid(db *gorm.DB, id int) structs.DiscsInfoRequest {
	var discussion structs.DiscsInfoRequest
	db.Table("Discussions").
		Select("Discussions.id, Discussions.title, Discussions.content, Discussions.created_at, Users.id,Users.username, Users.avatar").
		Joins("JOIN Users ON Discussions.Uid = Users.Uid").
		Where("Discussions.id = ?", id).First(&discussion)
	return discussion
}

// 添加讨论
func AddDiscussion(db *gorm.DB, title, content string, userId int) bool {
	if title == "" || content == "" {
		return false
	}
	err := db.Table("Discussions").
		Create(&tables.Discussion{UserId: userId, Title: title, Content: content, CreatedAt: time.Now()}).Error
	return err == nil
}

// 删除讨论
func DelDiscussion(db *gorm.DB, id int) bool {
	err := db.Table("Discussions").Where("id = ?", id).Delete(&tables.Discussion{}).Error
	return err == nil
}

// 添加评论
func AddComment(db *gorm.DB, content string, id, userId int, profanity bool) bool {
	return db.Table("Comments").Create(&tables.Comment{Id: id, UserId: userId, Content: content, CreatedAt: time.Now(), Profanity: profanity}).Error == nil
}

// 获取指定讨论ID的所有评论信息
func SelectCommentsByDid(db *gorm.DB, id int) []structs.CommentRequest {
	var comments []structs.CommentRequest
	db.Table("Comments").
		Select("Comments.id, Comments.discussion_id, Comments.content, Comments.created_at, Users.id,Users.username, Users.avatar,Comments.profanity").
		Joins("JOIN Users ON Comments.user_id = Users.id").
		Order("create_at desc").
		Where("Comments.discussion_id = ?", id).Find(&comments)
	return comments
}

// 删除指定评论
func DeleteCommentByCid(db *gorm.DB, id int) bool {
	return db.Table("Comments").Where("id = ?", id).Delete(&tables.Comment{}).Error == nil
}
