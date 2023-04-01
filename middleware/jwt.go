package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyx-0203/sgs-server/utils"
)

// JWT 基于JWT的认证中间件
func JWT() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Token放在Header的Authorization中，并使用Bearer开头
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			fmt.Println("请求头中auth为空")
			return
		}
		// 按空格分割
		// parts := strings.SplitN(authHeader, " ", 2)
		// if !(len(parts) == 2 && parts[0] == "Bearer") {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"code": 1,
		// 		"msg":  "请求头中auth格式有误",
		// 	})
		// 	c.Abort()
		// 	fmt.Println("请求头中auth格式有误")
		// 	fmt.Println(authHeader)
		// 	return
		// }
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "无效的Token",
			})
			c.Abort()
			fmt.Println("无效的Token")
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("UserID", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
