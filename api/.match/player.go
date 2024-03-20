package match

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/wyx-0203/sgs-server/models"
)

type Player struct {
	UserID    uint            `json:"id"`
	Position  uint            `json:"position"`
	Already   bool            `json:"already"`
	Name      string          `json:"nickname"`
	Character string          `json:"character"`
	conn      *websocket.Conn `json:"-"`
	Room      *Room           `json:"-"`
	send      chan []byte     `json:"-"`
}

func NewPlayer(personal *models.Personal, conn *websocket.Conn) *Player {
	p := &Player{
		UserID:    personal.UserID,
		Name:      personal.Name,
		Character: personal.Character,
		conn:      conn,
		send:      make(chan []byte),
	}

	AddPlayer <- p
	go p.read()
	go p.write()
	return p
}

func (p *Player) disconnect() {
	p.conn.Close()
	if p.Room != nil {
		p.Room.RemovePlayer(p)
	}
	RemovePlayer <- p
}

func (p *Player) read() {
	defer p.disconnect()
	for {
		// 读取消息
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		// 从websocket连接中读取的消息不做处理，直接广播
		p.Room.broadcast <- message
	}
}

func (p *Player) write() {
	defer p.disconnect()
	for message := range p.send {
		if err := p.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
			return
		}
	}
}
