package models

import (
	"log"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

func (c *Comment) SaveComment() (*Comment, error) {
	err := DB.Create(&c).Error
	if err != nil {
		log.Printf("SaveComment err,%s", err)
		return &Comment{}, err
	}
	log.Printf("SaveComment success,%+v", c)
	return c, nil
}

func SelectComments(postID int64) ([]Comment, error) {
	var comments []Comment
	if err := DB.Debug().Model(&Comment{}).Where("post_id=?", postID).Find(&comments).Error; err != nil {
		log.Printf("SelectComments err, %s", err)
		return comments, err
	}
	log.Printf("SelectComments size:%d", len(comments))
	return comments, nil
}
