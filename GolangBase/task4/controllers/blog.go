package controllers

import (
	"GolangBase/task4/models"
	"GolangBase/task4/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReqPost struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ReqUpdatePost struct {
	ReqPost
	ID uint `json:"id" binding:"required"`
}

type ReqComment struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"postID" binding:"required"`
}

func SavePost(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权的访问",
			"error":   err.Error(),
		})
		return
	}

	_, err = models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "用户未找到",
			"error":   err.Error(),
		})
		return
	}

	var req ReqPost
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "参数不正确",
			"error":   err.Error(),
		})
		return
	}
	p := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId,
	}

	if _, err := p.SavePost(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文章保存失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章保存成功",
		"data":    p.ID,
	})
}

func UpdatePost(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权的访问",
			"error":   err.Error(),
		})
		return
	}

	var req ReqUpdatePost
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "参数不正确",
			"error":   err.Error(),
		})
		return
	}
	var p models.Post
	p, err = models.GetPostByID(req.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "文章不存在, postID:" + strconv.Itoa(int(req.ID)),
			"error":   err.Error(),
		})
		return
	}
	if p.UserID != userId {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "不能修改其他用户的文章",
			"error":   err,
		})
		return
	}

	p.Title = req.Title
	p.Content = req.Content

	if _, err := p.UpdatePost(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文章修改失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章修改成功",
		"data":    p.ID,
	})
}

func SelectPostList(c *gin.Context) {
	posts, err := models.SelectPostList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文章列表查询失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章列表查询成功",
		"data":    posts,
	})
}

func DeletePost(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权的访问",
			"error":   err.Error(),
		})
		return
	}
	postIdStr := c.Query("postId")
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文章ID格式不正确",
			"error":   err.Error(),
		})
		return
	}
	var p models.Post
	p, err = models.GetPostByID(uint(postId))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "文章不存在, postID:" + strconv.Itoa(int(postId)),
			"error":   err.Error(),
		})
		return
	}
	if p.UserID != userId {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "不能删除其他用户的文章",
			"error":   err,
		})
		return
	}

	rowsAffected, err := p.DeletePost()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "文章删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
		"data":    rowsAffected,
	})
}

func SaveComment(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权的访问",
			"error":   err.Error(),
		})
		return
	}
	var req ReqComment
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err.Error(),
			"message": "登录失败",
		})
		return
	}
	comment := models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  userId,
	}

	if _, err := comment.SaveComment(); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err,
			"message": "评论保存失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "评论保存成功",
		"data":    comment.ID,
	})
}

func SelectComments(c *gin.Context) {
	postIdStr := c.Query("postId")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文章ID格式不正确",
			"error":   err.Error(),
		})
		return
	}
	comments, err := models.SelectComments(postId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err,
			"message": "评论查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "评论查询成功",
		"data":    comments,
	})
}
