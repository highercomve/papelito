package gamemodels

import (
	"time"

	"github.com/highercomve/papelito/modules/helpers/helpermodels"
)

const (
	ErrorGameNotCreated = "game can't be created"
)

type Player struct {
	helpermodels.Identification `json:",inline" bson:",inline"`

	Name  string  `json:"name" bson:"name"`
	Words []Words `json:"words" bson:"words"`
}

type Words struct {
	helpermodels.Identification `json:",inline" bson:",inline"`

	Title string `json:"title" bson:"title"`
}

type Scene struct {
	helpermodels.Identification `json:",inline" bson:",inline"`

	Number  int            `json:"number" bson:"number"`
	Words   []Words        `json:"words" bson:"words"`
	Results map[string]int `json:"results" bson:"results"`
}

type Group struct {
	helpermodels.Identification `json:",inline" bson:",inline"`

	Name    string            `json:"name" bson:"name"`
	Members map[string]Player `json:"members" bson:"members"`
	Words   map[string]map[string]bool
}

type Configuration struct {
	helpermodels.Identification `json:",inline" bson:",inline"`

	Groups          int           `json:"groups" form:"groups" validate:"required" bson:"groups"`
	Members         int           `json:"members" form:"members" validate:"required" bson:"members"`
	WordsForMembers int           `json:"words_for_members" form:"words_for_members" validate:"required" bson:"words_for_members"`
	Scenes          int           `json:"scenes" form:"scenes" bson:"scenes"`
	TurnDuration    time.Duration `json:"turn_duration" form:"turn_duration" validate:"required" bson:"turn_duration"`
}

type Game struct {
	helpermodels.Identification `json:",inline" bson:",inline"`
	helpermodels.Timestamp      `json:",inline" bson:",inline"`
	helpermodels.Ownership      `json:",inline" bson:",inline"`

	Players       map[string]Player `json:"players" bson:"players"`
	Words         []Words           `json:"words" bson:"words"`
	Groups        map[string]Group  `json:"groups" bson:"groups"`
	Scenes        []Scene           `json:"scenes" bson:"scenes"`
	Configuration *Configuration    `json:"configuration" bson:"configuration"`
}

func (game *Game) GetServicePrn() string {
	return game.ID
}
