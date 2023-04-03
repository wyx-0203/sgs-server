package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unrolled/secure"
	"github.com/wyx-0203/sgs-server/controllers"
	"github.com/wyx-0203/sgs-server/docs"
	"github.com/wyx-0203/sgs-server/global"
	"github.com/wyx-0203/sgs-server/match"
	"github.com/wyx-0203/sgs-server/models"
	"github.com/wyx-0203/sgs-server/utils"
)

func main() {
	models.InitDB()
	match.Init()

	r := gin.Default()
	router(r)

	if !global.SSL_IS_ON {
		r.Run(global.PORT) // 监听并在 0.0.0.0:8080 上启动服务
	} else {
		// https
		r.Use(TlsHandler())
		r.RunTLS(":"+global.PORT, global.SSL_CRT, global.SSL_CRT_KEY)
	}
}

func router(r *gin.Engine) {
	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 用户
	r.POST("/signup", controllers.SignUp)
	r.POST("/signin", controllers.SignIn)

	// 个人信息
	r.GET("/getUserInfo", controllers.GetUserInfo)

	r.GET("/websocket", controllers.ConnectWS)

	authed := r.Group("")
	authed.Use(utils.JWT())
	{
		// 游戏
		authed.GET("/createRoom", controllers.CreateRoom)
		authed.GET("/joinRoom", controllers.JoinRoom)
		authed.GET("/exitRoom", controllers.ExitRoom)
		authed.GET("/setAlready", controllers.SetAlready)
		authed.GET("/startGame", controllers.StartGame)

		// 个人信息
		authed.GET("/rename", controllers.Rename)
		authed.GET("/changeCharacter", controllers.ChangeCharacter)
	}
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
