package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wyx-0203/sgs-server/models"
)

func GetUserInfo(c *gin.Context) {
	// 查找个人信息
	id, _ := strconv.Atoi(c.Query("id"))
	p, err := models.FindPersonal(uint(id))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "success",
		"nickname":  p.Name,
		"character": p.Character,
		"win":       p.Win,
		"lose":      p.Lose,
	})
}

func Rename(c *gin.Context) {
	id := c.GetUint("UserID")
	name := c.Query("name")

	p, err := models.FindPersonal(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	models.Rename(p, name)

	player, ok := getPlayer(c)
	if ok {
		player.Name = name
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func ChangeCharacter(c *gin.Context) {
	id := c.GetUint("UserID")
	character := c.Query("character")

	p, err := models.FindPersonal(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	models.ChangeCharacter(p, character)

	player, ok := getPlayer(c)
	if ok {
		player.Character = character
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
