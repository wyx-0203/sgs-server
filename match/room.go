package match

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MODE_欢乐成双 = 0
	MODE_统帅双军 = 1
)

type Room struct {
	ID           int
	Players      []*Player
	Mode         int
	PlayerNumber int
	Owner        *Player
	inGame       bool
	broadcast    chan []byte
	mutex        sync.Mutex
}

var currentId = 1 //递增房间号

func NewRoom(mode int) *Room {
	r := &Room{
		Players:      make([]*Player, getPlayerNumber(mode)),
		Mode:         mode,
		PlayerNumber: getPlayerNumber(mode),
		broadcast:    make(chan []byte),
		mutex:        sync.Mutex{},
	}

	r.mutex.Lock()
	r.ID = currentId
	currentId = (currentId + 1) % 1000
	r.mutex.Unlock()

	AddRoom <- r
	go r.run()
	return r
}

func getPlayerNumber(mode int) int {
	if mode == MODE_欢乐成双 {
		return 4
	}
	return 2
}

func (r *Room) run() {
	for message := range r.broadcast {
		for _, p := range r.Players {
			if p != nil {
				p.send <- message
			}
		}
	}
}

func (r *Room) AddPlayer(p *Player) {
	// 添加玩家
	r.mutex.Lock()
	for i := range r.Players {
		if r.Players[i] == nil {
			r.Players[i] = p
			p.Position = uint(i)
			break
		}
	}

	p.Room = r
	p.Already = false
	if r.Owner == nil {
		r.Owner = p
	}

	fmt.Printf("room %d: ", r.ID)
	for _, i := range r.Players {
		if i != nil {
			fmt.Printf("%d", i.UserID)
		}
	}
	r.mutex.Unlock()

	// 发送消息
	message, _ := json.Marshal(joinRoomMsg{
		message: message{MsgType: "add_player"},
		Player:  *p,
	})
	r.broadcast <- message
}

func (r *Room) RemovePlayer(p *Player) {
	// 移除玩家
	r.mutex.Lock()
	r.Players[p.Position] = nil
	p.Room = nil

	rm := true
	for _, _p := range r.Players {
		if _p != nil {
			rm = false
			break
		}
	}
	if rm {
		RemoveRoom <- r
		r.mutex.Unlock()
		return
	}

	// if len(r.Players) == 0 {
	// 	// 移除房间
	// 	RemoveRoom <- r
	// 	r.mutex.Unlock()
	// 	return
	// } else
	if r.Owner == p {
		// 转让房主
		for _, _p := range r.Players {
			if _p != nil {
				r.Owner = _p
				break
			}
		}
	}
	r.mutex.Unlock()

	// 发送消息
	message, _ := json.Marshal(exitRoomMsg{
		message:  message{MsgType: "remove_player"},
		Position: p.Position,
		OwnerPos: r.Owner.Position,
	})
	r.broadcast <- message
}

func (r *Room) SetAlready(p *Player) {
	p.Already = !p.Already

	message, _ := json.Marshal(setAlreadyMsg{
		message:  message{MsgType: "set_already"},
		Position: p.Position,
		Already:  p.Already,
	})
	p.Room.broadcast <- message
}

func (r *Room) StartGame() {
	r.inGame = true

	p := []uint{}
	for _, i := range r.Players {
		p = append(p, i.UserID)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})

	message, _ := json.Marshal(startGameMsg{
		message: message{MsgType: "start_game"},
		Players: p,
	})
	r.broadcast <- message

	// 取消准备
	for _, i := range r.Players {
		r.SetAlready(i)
	}
}
