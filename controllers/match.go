package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wyx-0203/sgs-server/match"
	"github.com/wyx-0203/sgs-server/models"
)

var upgrader = websocket.Upgrader{
	// 防止跨域站点伪造请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建player对象，浏览器websocket连接不支持自定义标头，所以将UserID放在参数里
	id, _ := strconv.Atoi(c.Query("user_id"))
	p, _ := models.FindPersonal(uint(id))
	match.NewPlayer(p, conn)
}

// 通过token中的id获取player对象
func getPlayer(c *gin.Context) (*match.Player, bool) {
	id := c.GetUint("UserID")
	p, ok := match.Players[id]
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "未登录",
		})
	}
	return p, ok
}

// 快速加入
func JoinRoom(c *gin.Context) {
	// 查找玩家
	p, ok := getPlayer(c)
	if !ok {
		return
	}

	// 查找房间
	mode, _ := strconv.Atoi(c.Query("mode"))
	r := match.QuickFind(mode)

	// 创建房间
	if r == nil {
		r = match.NewRoom(mode)
	}

	// 进入房间
	r.AddPlayer(p)

	fmt.Printf("join room room_id:%d user_id:%d\n", r.ID, p.UserID)
	fmt.Printf("\n")

	// 返回房间信息
	getRoomInfo(c, r)
}

// 退出房间
func ExitRoom(c *gin.Context) {
	// 查找玩家
	p, ok := getPlayer(c)
	if !ok {
		return
	}

	p.Room.RemovePlayer(p)

	// 返回
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出房间",
	})
}

// 创建房间
func CreateRoom(c *gin.Context) {
	// 查找玩家
	p, ok := getPlayer(c)
	if !ok {
		return
	}

	// 创建房间
	mode, _ := strconv.Atoi(c.Query("mode"))
	r := match.NewRoom(mode)

	// 进入房间
	r.AddPlayer(p)

	// 返回
	getRoomInfo(c, r)
}

// 获取房间信息
func getRoomInfo(c *gin.Context, r *match.Room) {
	players := []*match.Player{}
	for _, p := range r.Players {
		if p != nil {
			players = append(players, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"mode":      r.Mode,
		"players":   players,
		"owner_pos": r.Owner.Position,
	})
}

// func GetRoomInfo(c *gin.Context) {
// 	// 查找房间
// 	id, _ := strconv.Atoi(c.Query("room_id"))
// 	r, ok := match.Rooms[id]
// 	if !ok {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    1,
// 			"message": "房间不存在",
// 		})
// 		return
// 	}
// 	getRoomInfo(c, r)
// }

func SetAlready(c *gin.Context) {
	p, ok := getPlayer(c)
	if !ok {
		return
	}

	p.Room.SetAlready(p)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 开始游戏
func StartGame(c *gin.Context) {
	p, ok := getPlayer(c)
	if !ok {
		return
	}

	p.Room.StartGame()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
