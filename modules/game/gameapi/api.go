package gameapi

import (
	"fmt"
	"sync"
	"time"

	"github.com/highercomve/papelito/modules/game/gameservice"
	"github.com/labstack/echo/v4"
)

var gameService *gameservice.GameMachine

// Load Create new auth service
func Load(e *echo.Group) *echo.Group {
	g := e.Group("/games")
	gameService = gameservice.NewGameMachine()

	g.GET("", GetCreateGame)
	g.POST("", CreateGame)
	g.GET("/:id", GetGame)
	g.PUT("/:id", UpdateGame)
	g.GET("/:id/sse", handleSSE)

	fmt.Printf("%++v", g)
	return g
}

var (
	players = make(map[string]map[chan string]bool)
	mutex   = &sync.Mutex{}
)

func handleSSE(c echo.Context) error {
	id := c.Param("id")

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	// Create a channel for this client
	messageChan := make(chan string)
	defer close(messageChan)

	mutex.Lock()
	if _, ok := players[id]; !ok {
		players[id] = make(map[chan string]bool)
	}
	players[id][messageChan] = true
	mutex.Unlock()

	// Send an initial event to acknowledge the connection
	fmt.Fprintf(c.Response(), "data: Connected to game %s\n\n", id)
	c.Response().Flush()

	// Periodically send a heartbeat to keep the connection alive
	go func() {
		for {
			select {
			case <-c.Request().Context().Done():
				// Client disconnected
				mutex.Lock()
				delete(players[id], messageChan)
				if len(players[id]) == 0 {
					delete(players, id)
				}
				mutex.Unlock()
				return

			case message := <-messageChan:
				// Broadcast the message to all players in the game
				broadcast(id, message)

			case <-time.After(30 * time.Second):
				// Send a heartbeat to keep the connection alive
				fmt.Fprintf(c.Response(), ":heartbeat\n\n")
				c.Response().Flush()
			}
		}
	}()

	return nil
}

func broadcast(id, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameplayers, ok := players[id]; ok {
		for client := range gameplayers {
			select {
			case client <- message:
			default:
				// Remove the client if the channel is full
				delete(gameplayers, client)
			}
		}
	}
}
