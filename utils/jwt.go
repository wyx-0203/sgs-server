package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wyx-0203/sgs-server/global"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func GenerateToken(id uint) (string, error) {
	// 创建 Claims
	claims := &Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}
	// 生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	return token.SignedString([]byte(global.JWT_SECRET))
}

func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

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
		mc, err := ParseToken(token)
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
