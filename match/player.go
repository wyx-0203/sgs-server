package match

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/wyx-0203/sgs-server/models"
)

// const (
// 	// Max wait time when writing message to peer
// 	writeWait = 5 * time.Second

// 	// Max time till next pong from peer
// 	pongWait = 20 * time.Second

// 	// Send ping interval, must be less then pong wait time
// 	pingPeriod = (pongWait * 9) / 10
// )

type Player struct {
	UserID    uint   `json:"id"`
	Position  uint   `json:"position"`
	Already   bool   `json:"already"`
	Name      string `json:"nickname"`
	Character string `json:"character"`
	Room      *Room  `json:"-"`

	conn *websocket.Conn `json:"-"`
	send chan []byte     `json:"-"`
}

func NewPlayer(personal *models.Personal, conn *websocket.Conn) *Player {
	p := &Player{
		UserID:    personal.UserID,
		Name:      personal.Name,
		Character: personal.Character,

		conn: conn,
		send: make(chan []byte),
	}

	AddPlayer <- p
	go p.read()
	go p.write()
	// go p.ping()
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
	// p.conn.SetReadDeadline(time.Now().Add(pongWait))
	// p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		// 读取消息
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("unexpected close error: %v", err)
			}
			break
		}

		// 从websocket连接中读取的消息不做处理，直接广播
		p.Room.broadcast <- message
	}
}

func (p *Player) write() {
	for message := range p.send {
		if err := p.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// func (p *Player) ping() {
// 	defer p.disconnect()

// 	ticker := time.NewTicker(pingPeriod)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		p.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 		if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 			fmt.Println("player " + p.Name + ": ping pong error")
// 			return
// 		}
// 	}
// }
