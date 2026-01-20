package handler

import (
	"github.com/Seeker32/my-blog/internal/request"
	"github.com/Seeker32/my-blog/internal/response"
	"github.com/Seeker32/my-blog/internal/service"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	user, err := service.CreateUser(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}, "注册成功")
}
