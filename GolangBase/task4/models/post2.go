package models

import (
	"log"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint
	User     User
	Comments []Comment `gorm:"foreignKey:PostID"`
}

func (p *Post) SavePost() (*Post, error) {
	err := DB.Create(&p).Error
	if err != nil {
		log.Printf("savePost err,%s", err)
		return &Post{}, err
	}
	log.Printf("savePost success,%+v", p)
	return p, nil
}

func (p *Post) UpdatePost() (*Post, error) {
	err := DB.Debug().Model(&p).Updates(&p).Error
	if err != nil {
		log.Printf("UpdatePost err,%s", err)
		return &Post{}, err
	}
	log.Printf("UpdatePost success,%+v", p)
	return p, nil
}

func (p *Post) DeletePost() (int64, error) {
	result := DB.Debug().Delete(&p)
	if result.Error != nil {
		log.Printf("DeletePost err,%s", result.Error)
		return 0, result.Error
	}
	log.Printf("DeletePost success,%d", result.RowsAffected)
	return result.RowsAffected, nil
}

func GetPostByID(id uint) (Post, error) {
	var p Post
	if err := DB.Debug().First(&p, id).Error; err != nil {
		log.Printf("GetPostByID err,%s", err)
		return p, err
	}
	log.Printf("GetPostByID success,%+v", p)
	return p, nil
}

func SelectPostList() ([]Post, error) {
	var posts []Post
	if err := DB.Debug().Model(&Post{}).Preload("Comments").Find(&posts).Error; err != nil {
		log.Printf("SelectPostList err, %s", err)
		return posts, err
	}
	log.Printf("SelectPostList size:%d", len(posts))
	return posts, nil
}
