package game

//nolint
type GameForm struct {
	Params struct {
		PlayerName string `json:"playerName"`
		GameName   string `json:"gameName"`
		Spies      int    `json:"spyNumber"`
		MaxNumbers int    `json:"maxNumbers"`
		IsPrivate  bool   `json:"isPrivate"`
	} `json:"params"`
}

type Round struct {
	Number  int      `json:"number"`
	Word    string   `json:"word"`
	SpyWord string   `json:"spyWord"`
	Spies   []string `json:"spies"`
	Winner  string   `json:"winner"`
}

type Game struct {
	RoomID      string   `json:"roomID"`
	Players     []string `json:"players"`
	Word        string   `json:"word"`
	SpyWord     string   `json:"spyWord"`
	Spies       []string `json:"spies"`
	RoomName    string   `json:"roomName"`
	IsPublic    bool     `json:"isPublic"`
	MaxPlayers  int      `json:"maxPlayers"`
	GameStarted bool     `json:"gameStarted"`
	Rounds      []Round  `json:"rounds"`
	CreatedAt   int64    `json:"createdAt"`
}

