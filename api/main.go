package main

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wyx-0203/sgs-server/controllers"
	"github.com/wyx-0203/sgs-server/docs"

	// "github.com/wyx-0203/sgs-server/match"
	"github.com/wyx-0203/sgs-server/middleware"
	"github.com/wyx-0203/sgs-server/models"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest"
)

func main() {
	// 使用 InClusterConfig 从 Pod 内部获取 Kubernetes 集群配置
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	// 处理错误
	// }

	// // 创建 Kubernetes 客户端
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	// 处理错误
	// }

	// // 获取 Core V1 API 的接口
	// coreV1API := clientset.CoreV1()

	// // 获取 ConfigMap 数据
	// configMap, err := coreV1API.ConfigMaps("namespace").Get(context.TODO(), "configmap-name", metav1.GetOptions{})
	// if err != nil {
	// 	// 处理错误
	// }

	// 处理获取到的 ConfigMap 数据

	models.InitDB()
	controllers.InitGrpc(os.Getenv("ROOM_SERVICE"))
	// match.Init()

	// Set up a connection to the server.
	// log.Println("aaa")
	// resp, err := c.CreateRoom(context.Background(), &room.CreateRoomRequest{UserId: 1})
	// if err != nil {
	// 	log.Fatalf("111 %v", err)
	// }
	// println(resp.RoomId)

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

	// r.GET("/websocket", controllers.ConnectWS)

	authed := r.Group("")
	authed.Use(middleware.JWT())
	{
		// 游戏
		authed.GET("/createRoom", controllers.CreateRoom)
		authed.GET("/joinRoom", controllers.JoinRoom)
		// authed.GET("/exitRoom", controllers.ExitRoom)
		// authed.GET("/getRoomInfo", controllers.GetRoomInfo)
		// authed.GET("/setAlready", controllers.SetAlready)
		// authed.GET("/startGame", controllers.StartGame)

		// 个人信息
		authed.GET("/rename", controllers.Rename)
		authed.GET("/changeCharacter", controllers.ChangeCharacter)
	}
}
