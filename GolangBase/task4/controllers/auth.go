package controllers

import (
	"GolangBase/task4/models"
	"GolangBase/task4/utils/token"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ReqRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ReqLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func Register(c *gin.Context) {
	var req ReqRegister
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	u := models.User{
		Username: req.Username,
		Password: req.Password,
	}
	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "注册失败",
			"err":     err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功了",
		"data":    req,
	})
}

func Login(c *gin.Context) {
	var req ReqLogin
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err.Error(),
			"message": "登录失败",
		})
		return
	}
	u := models.User{
		Username: req.Username,
		Password: req.Password,
	}
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err,
			"message": "用户名或密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

func CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权的访问",
			"error":   err.Error(),
		})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "用户未找到",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户成功",
		"data":    user,
	})
}
