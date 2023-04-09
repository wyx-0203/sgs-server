package match

import "fmt"

var (
	Rooms   = map[int]*Room{}    // 所有房间
	Players = map[uint]*Player{} //所有在线玩家

	AddRoom      = make(chan *Room)
	RemoveRoom   = make(chan *Room)
	AddPlayer    = make(chan *Player)
	RemovePlayer = make(chan *Player)
)

func Init() {
	go handleRoom()
	go handlePlayer()
}

func handleRoom() {
	for {
		select {
		case r := <-AddRoom:
			Rooms[r.ID] = r
		case r := <-RemoveRoom:
			delete(Rooms, r.ID)
		}

		fmt.Printf("rooms: ")
		for i := range Rooms {
			fmt.Printf("%d ", i)
		}
		fmt.Printf("\n")
	}
}

func handlePlayer() {
	for {
		select {
		case p := <-AddPlayer:
			Players[p.UserID] = p
		case p := <-RemovePlayer:
			delete(Players, p.UserID)
		}

		fmt.Printf("players: ")
		for i := range Players {
			fmt.Printf("%d ", i)
		}
		fmt.Printf("\n")
	}
}

func QuickFind(mode int) *Room {
	// 找到第一个还有空位的房间
	for _, r := range Rooms {
		if r.inGame || r.Mode != mode {
			continue
		}

		for _, p := range r.Players {
			if p == nil {
				return r
			}
		}
	}
	return nil
}
