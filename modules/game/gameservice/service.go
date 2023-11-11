package gameservice

import (
	"context"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/highercomve/papelito/modules/game/gamemodels"
	"github.com/highercomve/papelito/modules/game/gamerepo"
	"github.com/highercomve/papelito/modules/game/statemachine"
	"github.com/highercomve/papelito/modules/helpers/helpermodels"
	"github.com/highercomve/papelito/modules/helpers/helperrepo"
)

const (
	CreateGameEvent    = statemachine.EventType("create-game")
	OnCreatedGameEvent = statemachine.EventType("on-created-game")
)

type GameMachine struct {
	machine statemachine.IStateMachine
	storage *gamerepo.Repo

	Games map[string]gamemodels.Game
}

func NewGameMachine() *GameMachine {
	statemachine := statemachine.NewGameMachine()
	repo, err := gamerepo.GetRepo()
	if err != nil {
		return nil
	}
	machine := &GameMachine{
		machine: statemachine,
		storage: repo,
	}

	return machine
}

func (m *GameMachine) CreateGame(ctx context.Context, config *gamemodels.Configuration) (*gamemodels.Game, error) {
	game := &gamemodels.Game{
		Identification: helpermodels.NewIdentification("game"),
		Timestamp:      helpermodels.NewTimeStamp(),
		Ownership:      helpermodels.Ownership{},
		Configuration:  config,
		Players:        map[string]gamemodels.Player{},
		Words:          make([]gamemodels.Words, config.WordsForMembers*config.Members),
		Groups:         map[string]gamemodels.Group{},
	}
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	for i := 0; i < game.Configuration.Groups; i++ {
		group := gamemodels.Group{
			Identification: helpermodels.NewIdentification("group"),
			Name:           nameGenerator.Generate(),
			Members:        map[string]gamemodels.Player{},
			Words:          map[string]map[string]bool{},
		}
		game.Groups[group.ID] = group
	}

	err := m.storage.Repo.UpdateOne(ctx, game, true)
	return game, err
}

func (m *GameMachine) GetGame(ctx context.Context, id string, p helperrepo.Map) (game *gamemodels.Game, err error) {
	game = &gamemodels.Game{}
	err = m.storage.Repo.FindByID(ctx, id, p, game)
	return game, err
}
