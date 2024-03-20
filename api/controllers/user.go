package controllers

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyx-0203/sgs-server/models"
	"github.com/wyx-0203/sgs-server/utils"
)

func md5Encode(str string) string {
	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

// @Summary		注册
// @Consume		multipart/form-data
// @param		username	formData	string	true	"username"
// @param		password	formData	string	true	"password"
// @Router		/signup [post]
func SignUp(c *gin.Context) {
	// 用户名已存在
	if _, err := models.FindUserByName(c.Request.FormValue("username")); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名已存在",
		})
		return
	}

	// 创建用户
	user := &models.User{
		Username: c.Request.FormValue("username"),
		Password: md5Encode(c.Request.FormValue("password")),
	}
	models.CreateUser(user)
	models.CreatePersonal(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
	})
}

// @Summary		登录
// @Consume		multipart/form-data
// @param		username	formData	string	true	"username"
// @param		password	formData	string	true	"password"
// @Router		/signin [post]
func SignIn(c *gin.Context) {
	// 验证用户名与密码
	username := c.Request.FormValue("username")
	password := md5Encode(c.Request.FormValue("password"))
	user, err := models.FindUserByNameAndPwd(username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名或密码错误",
		})
		return
	}

	// 登录成功，返回token
	token, _ := utils.GenerateToken(user.ID)
	fmt.Println(token)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"token":   token,
		"user_id": user.ID,
	})
}
