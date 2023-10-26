package gameservice

import (
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/highercomve/papelito/modules/game/statemachine"
	"github.com/highercomve/papelito/modules/helper/helpermodels"
)

const (
	CreateGameEvent    = statemachine.EventType("create-game")
	OnCreatedGameEvent = statemachine.EventType("on-created-game")
)

const (
	ErrorGameNotCreated = "game can't be created"
)

type Player struct {
	helpermodels.Identification `json:",inline"`

	Name  string  `json:"name"`
	Words []Words `json:"words"`
}

type Words struct {
	helpermodels.Identification `json:",inline"`

	Title string `json:"title"`
}

type Scene struct {
	helpermodels.Identification `json:",inline"`

	Number  int            `json:"number"`
	Words   []Words        `json:"words"`
	Results map[string]int `json:"results"`
}

type Group struct {
	helpermodels.Identification `json:",inline"`

	Name    string            `json:"name"`
	Members map[string]Player `json:"members"`
	Words   map[string]map[string]bool
}

type Configuration struct {
	helpermodels.Identification `json:",inline"`

	Groups          int           `json:"groups" form:"groups" validate:"required"`
	Members         int           `json:"members" form:"members" validate:"required"`
	WordsForMembers int           `json:"words_for_members" form:"words_for_members" validate:"required"`
	Scenes          int           `json:"scenes" form:"scenes"`
	TurnDuration    time.Duration `json:"turn_duration" form:"turn_duration" validate:"required"`
}

type Game struct {
	helpermodels.Identification
	helpermodels.Timestamp

	Players       map[string]Player `json:"players"`
	Words         []Words           `json:"words"`
	Groups        map[string]Group  `json:"groups"`
	Scenes        []Scene           `json:"scenes"`
	Configuration *Configuration    `json:"configuration"`
}

type GameMachine struct {
	machine statemachine.IStateMachine

	Games map[string]Game
}

func NewGameMachine() *GameMachine {
	statemachine := statemachine.NewGameMachine()
	machine := &GameMachine{
		machine: statemachine,
		Games:   map[string]Game{},
	}

	return machine
}

func (m *GameMachine) CreateGame(config *Configuration) (*Game, error) {
	game := Game{
		Identification: helpermodels.NewIdentification("game"),
		Timestamp:      helpermodels.NewTimeStamp(),
		Configuration:  config,
		Players:        map[string]Player{},
		Words:          make([]Words, config.WordsForMembers*config.Members),
		Groups:         map[string]Group{},
	}
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	for i := 0; i < game.Configuration.Groups; i++ {
		group := Group{
			Identification: helpermodels.NewIdentification("group"),
			Name:           nameGenerator.Generate(),
			Members:        map[string]Player{},
			Words:          map[string]map[string]bool{},
		}
		game.Groups[group.ID] = group
	}

	m.Games[game.ID] = game

	return &game, nil
}
