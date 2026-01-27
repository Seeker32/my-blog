package handler

import (
	"github.com/Seeker32/my-blog/internal/request"
	"github.com/Seeker32/my-blog/internal/response"
	"github.com/Seeker32/my-blog/internal/service"
	"github.com/gin-gonic/gin"
)

// Register 用户注册处理函数
//
//	@Summary		用户注册
//	@Description	用户注册
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.UserRegisterRequest	true	"用户注册信息"
//	@Success		200		{object}	user.UserRegisterResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
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
