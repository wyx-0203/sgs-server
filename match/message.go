package match

type message struct {
	MsgType string `json:"msg_type"`
}

type joinRoomMsg struct {
	message
	Player Player `json:"player"`
}

type exitRoomMsg struct {
	message
	Position uint `json:"position"`
	OwnerPos uint `json:"owner_pos"`
}

type setAlreadyMsg struct {
	message
	Position uint `json:"position"`
	Already  bool `json:"already"`
}

type startGameMsg struct {
	message
	Players []uint `json:"players"`
}
