package battlesnake

const (
	MOVE_UP    = "up"
	MOVE_RIGHT = "right"
	MOVE_DOWN  = "down"
	MOVE_LEFT  = "left"
)

var Moves = []string{MOVE_UP, MOVE_RIGHT, MOVE_DOWN, MOVE_LEFT}

var MovesMap = map[string]Point{
	MOVE_UP:    {0, 1},
	MOVE_RIGHT: {1, 0},
	MOVE_DOWN:  {0, -1},
	MOVE_LEFT:  {-1, 0},
}

type Ruleset struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Settings struct {
		FoodSpawnChance     int    `json:"foodSpawnChance"`
		MinimumFood         int    `json:"minimumFood"`
		HazardDamagePerTurn int    `json:"hazardDamagePerTurn"`
		HazardMap           string `json:"hazardMap"`
		HazardMapAuthor     string `json:"hazardMapAuthor"`
		Royale              struct {
			ShrinkEveryNTurns int `json:"shrinkEveryNTurns"`
		} `json:"royale"`
		Squad struct {
			AllowBodyCollisions bool `json:"allowBodyCollisions"`
			SharedElimination   bool `json:"sharedElimination"`
			SharedHealth        bool `json:"sharedHealth"`
			SharedLength        bool `json:"sharedLength"`
		} `json:"squad"`
	} `json:"settings"`
}

type Snake struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Latency        string  `json:"latency"`
	Health         int     `json:"health"`
	Body           []Point `json:"body"`
	Head           Point   `json:"head"`
	Length         int     `json:"length"`
	Shout          string  `json:"shout"`
	Squad          string  `json:"squad"`
	Customizations struct {
		Color string `json:"color"`
		Head  string `json:"head"`
		Tail  string `json:"tail"`
	} `json:"customizations"`
}

type Board struct {
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Snakes  []Snake       `json:"snakes"`
	Food    []Point       `json:"food"`
	Hazards []interface{} `json:"hazards"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type State struct {
	Game struct {
		ID      string  `json:"id"`
		Ruleset Ruleset `json:"ruleset"`
		Map     string  `json:"map"`
		Timeout int     `json:"timeout"`
		Source  string  `json:"source"`
	} `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}
