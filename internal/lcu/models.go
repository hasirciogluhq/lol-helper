package lcu

// Summoner oyuncu bilgisi
type Summoner struct {
	AccountID     int64  `json:"accountId"`
	DisplayName   string `json:"displayName"`
	GameName      string `json:"gameName"`
	InternalName  string `json:"internalName"`
	ProfileIconID int    `json:"profileIconId"`
	Puuid         string `json:"puuid"`
	SummonerID    int64  `json:"summonerId"`
	SummonerLevel int    `json:"summonerLevel"`
}

// GameFlowSession oyun akış durumu
type GameFlowSession struct {
	Phase string `json:"phase"` // None, Lobby, Matchmaking, ChampSelect, InProgress, etc.
}

// GameData oyun verisi
type GameData struct {
	Phase       string              `json:"phase"`
	GameTime    int                 `json:"gameTime"`
	ChampSelect *ChampSelectSession `json:"champSelect,omitempty"`
	InGame      *InGameInfo         `json:"inGame,omitempty"`
}

// ChampSelectSession champion seçim bilgisi
type ChampSelectSession struct {
	Actions       [][]ChampSelectAction `json:"actions"`
	AlliedTeam    []ChampSelectPlayer   `json:"myTeam"`
	EnemyTeam     []ChampSelectPlayer   `json:"theirTeam"`
	LocalPlayerID int64                 `json:"localPlayerCellId"`
	Timer         ChampSelectTimer      `json:"timer"`
}

// ChampSelectAction champion seçim aksiyonu
type ChampSelectAction struct {
	ActorCellID int64  `json:"actorCellId"`
	ChampionID  int    `json:"championId"`
	Completed   bool   `json:"completed"`
	ID          int64  `json:"id"`
	Type        string `json:"type"` // pick, ban
}

// ChampSelectPlayer champion select'teki oyuncu
type ChampSelectPlayer struct {
	CellID       int64  `json:"cellId"`
	ChampionID   int    `json:"championId"`
	ChampionName string `json:"championPickIntent"`
	SummonerID   int64  `json:"summonerId"`
	Team         int    `json:"team"`
}

// ChampSelectTimer sayaç bilgisi
type ChampSelectTimer struct {
	AdjustedTimeLeftInPhase int64  `json:"adjustedTimeLeftInPhase"`
	InternalNowInEpochMs    int64  `json:"internalNowInEpochMs"`
	IsInfinite              bool   `json:"isInfinite"`
	Phase                   string `json:"phase"`
	TotalTimeInPhase        int64  `json:"totalTimeInPhase"`
}

// InGameInfo oyun içi bilgi
type InGameInfo struct {
	GameTime int      `json:"gameTime"`
	Players  []Player `json:"players"`
}

// Player oyuncu bilgisi
type Player struct {
	SummonerName string `json:"summonerName"`
	ChampionName string `json:"championName"`
	ChampionID   int    `json:"championId"`
	Team         string `json:"team"` // ORDER, CHAOS
	Position     string `json:"position"`
	Level        int    `json:"level"`
	Items        []Item `json:"items"`
	Gold         int    `json:"gold"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
}

// Item item bilgisi
type Item struct {
	ItemID      int    `json:"itemID"`
	Count       int    `json:"count"`
	DisplayName string `json:"displayName"`
	Slot        int    `json:"slot"`
}
