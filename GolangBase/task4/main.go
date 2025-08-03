package main

import (
	"GolangBase/task4/controllers"
	"GolangBase/task4/middlewares"
	"GolangBase/task4/models"
	"github.com/gin-gonic/gin"
)

func init() {
	models.ConnectDB()
}

func main() {
	r := gin.Default()
	admin := r.Group("/auth")
	{
		admin.POST("/register", controllers.Register)
		admin.POST("/login", controllers.Login)
		admin.GET("/user", middlewares.JWtAuthMiddleware(), controllers.CurrentUser)
	}

	blog := r.Group("/blog")
	{
		post := blog.Group("/post")
		{
			post.POST("/create", middlewares.JWtAuthMiddleware(), controllers.SavePost)
			post.POST("/update", middlewares.JWtAuthMiddleware(), controllers.UpdatePost)
			post.DELETE("/delete", middlewares.JWtAuthMiddleware(), controllers.DeletePost)
			post.GET("/read", controllers.SelectPostList)
		}
		comment := blog.Group("/comment")
		{
			comment.POST("/create", middlewares.JWtAuthMiddleware(), controllers.SaveComment)
			comment.GET("/read", controllers.SelectComments)
		}
	}

	r.Run(":8080")
}
