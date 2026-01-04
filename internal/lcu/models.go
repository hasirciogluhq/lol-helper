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

// LiveGameData Live Client Data API yanıtı
type LiveGameData struct {
	ActivePlayer LiveActivePlayer `json:"activePlayer"`
	AllPlayers   []LivePlayer     `json:"allPlayers"`
	GameData     LiveGameStats    `json:"gameData"`
}

type LiveActivePlayer struct {
	SummonerName string `json:"summonerName"`
	Level        int    `json:"level"`
	CurrentGold  float64 `json:"currentGold"`
}

type LivePlayer struct {
	ChampionName    string       `json:"championName"`
	IsBot           bool         `json:"isBot"`
	IsDead          bool         `json:"isDead"`
	Items           []LiveItem   `json:"items"`
	Level           int          `json:"level"`
	Position        string       `json:"position"`
	RawChampionName string       `json:"rawChampionName"`
	RespawnTimer    float64      `json:"respawnTimer"`
	Runes           LiveRunes    `json:"runes"`
	Scores          LiveScores   `json:"scores"`
	SkinID          int          `json:"skinID"`
	SummonerName    string       `json:"summonerName"`
	SummonerSpells  LiveSpells   `json:"summonerSpells"`
	Team            string       `json:"team"` // ORDER, CHAOS
}

type LiveItem struct {
	CanUse      bool   `json:"canUse"`
	Consumable  bool   `json:"consumable"`
	Count       int    `json:"count"`
	DisplayName string `json:"displayName"`
	ItemID      int    `json:"itemID"`
	Price       int    `json:"price"`
	RawDescription string `json:"rawDescription"`
	RawDisplayName string `json:"rawDisplayName"`
	Slot        int    `json:"slot"`
}

type LiveRunes struct {
	Keystone LiveRune `json:"keystone"`
	PrimaryRuneTree LiveRuneTree `json:"primaryRuneTree"`
	SecondaryRuneTree LiveRuneTree `json:"secondaryRuneTree"`
}

type LiveRune struct {
	DisplayName string `json:"displayName"`
	ID          int    `json:"id"`
	RawDescription string `json:"rawDescription"`
}

type LiveRuneTree struct {
	DisplayName string `json:"displayName"`
	ID          int    `json:"id"`
	RawDescription string `json:"rawDescription"`
}

type LiveScores struct {
	Assists    int `json:"assists"`
	CreepScore int `json:"creepScore"`
	Deaths     int `json:"deaths"`
	Kills      int `json:"kills"`
	WardScore  float64 `json:"wardScore"`
}

type LiveSpells struct {
	SummonerSpellOne LiveSpell `json:"summonerSpellOne"`
	SummonerSpellTwo LiveSpell `json:"summonerSpellTwo"`
}

type LiveSpell struct {
	DisplayName string `json:"displayName"`
	RawDescription string `json:"rawDescription"`
	RawDisplayName string `json:"rawDisplayName"`
}

type LiveGameStats struct {
	GameTime float64 `json:"gameTime"`
	MapName  string  `json:"mapName"`
	MapNumber int    `json:"mapNumber"`
	MapTerrain string `json:"mapTerrain"`
}
