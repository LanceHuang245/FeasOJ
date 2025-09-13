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

	db.Table("discussions").
		Select("discussions.id, discussions.title, users.username, discussions.created_at").
		Joins("JOIN users ON discussions.user_id = users.id").
		Order("discussions.created_at desc").Count(&total).
		Offset((page - 1) * itemsPerPage).Limit(itemsPerPage).Find(&discussRequests)

	return discussRequests, int(total)
}

// 获取指定id讨论及User表中发帖人的头像
func SelectDiscussionByDid(db *gorm.DB, id int) structs.DiscsInfoRequest {
	var discussion structs.DiscsInfoRequest
	db.Table("discussions").
		Select("discussions.id, discussions.title, discussions.content, discussions.created_at, users.id,users.username, users.avatar").
		Joins("JOIN users ON discussions.user_id = users.id").
		Where("discussions.id = ?", id).First(&discussion)
	return discussion
}

// 添加讨论
func AddDiscussion(db *gorm.DB, title, content, userId string) bool {
	if title == "" || content == "" {
		return false
	}
	err := db.Table("discussions").
		Create(&tables.Discussions{UserId: userId, Title: title, Content: content, CreatedAt: time.Now()}).Error
	return err == nil
}

// 删除讨论
func DelDiscussion(db *gorm.DB, id int) bool {
	err := db.Table("discussions").Where("id = ?", id).Delete(&tables.Discussions{}).Error
	if err != nil {
		return false
	}
	err = db.Table("comments").Where("discussion_id = ?", id).Delete(&tables.Comments{}).Error
	if err != nil {
		return false
	}
	return true
}

// 添加评论
func AddComment(db *gorm.DB, content, userId string, DiscussionId int, profanity bool) bool {
	return db.Table("comments").Create(&tables.Comments{DiscussionId: DiscussionId, UserId: userId, Content: content, CreatedAt: time.Now(), Profanity: profanity}).Error == nil
}

// 获取指定讨论ID的所有评论信息
func SelectCommentsByDid(db *gorm.DB, id int) []structs.CommentRequest {
	var comments []structs.CommentRequest
	db.Table("comments").
		Select("comments.id as id, comments.discussion_id, comments.content, comments.created_at, users.id as user_id, users.username, users.avatar, comments.profanity").
		Joins("JOIN users ON comments.user_id = users.id").
		Order("created_at desc").
		Where("comments.discussion_id = ?", id).Find(&comments)
	return comments
}

// 删除指定评论
func DeleteCommentByCid(db *gorm.DB, id int) bool {
	return db.Table("comments").Where("id = ?", id).Delete(&tables.Comments{}).Error == nil
}
