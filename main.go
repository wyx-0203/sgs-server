package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wyx-0203/sgs-server/controllers"
	"github.com/wyx-0203/sgs-server/docs"
	"github.com/wyx-0203/sgs-server/match"
	"github.com/wyx-0203/sgs-server/middleware"
	"github.com/wyx-0203/sgs-server/models"
)

func main() {
	models.InitDB()
	match.Init()

	r := gin.Default()
	router(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func router(r *gin.Engine) {
	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 静态资源
	// r.Static("/assets", "./assets")
	// r.LoadHTMLFiles("templates/index.html")
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{})
	// })

	// 用户
	r.POST("/signup", controllers.SignUp)
	r.POST("/signin", controllers.SignIn)

	// 个人信息
	r.GET("/getUserInfo", controllers.GetUserInfo)

	r.GET("/websocket", controllers.ConnectWS)

	authed := r.Group("")
	authed.Use(middleware.JWT())
	{
		// 游戏
		authed.GET("/createRoom", controllers.CreateRoom)
		authed.GET("/joinRoom", controllers.JoinRoom)
		authed.GET("/exitRoom", controllers.ExitRoom)
		// authed.GET("/getRoomInfo", controllers.GetRoomInfo)
		authed.GET("/setAlready", controllers.SetAlready)
		authed.GET("/startGame", controllers.StartGame)

		// 个人信息
		authed.GET("/rename", controllers.Rename)
		authed.GET("/changeCharacter", controllers.ChangeCharacter)
	}
}
