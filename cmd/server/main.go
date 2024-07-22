package main

import (
	"log"
	"websocket/internal"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	room := internal.NewRoom()
	go room.Update()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			// if string(c.Request().Header.Peek("X-Api-Key")) != "123" {
			// 	return fiber.ErrBadRequest
			// }

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:name", websocket.New(func(conn *websocket.Conn) {

		var (
			p     = room.PlayerAdd(conn)
			err   error
			event internal.Event
		)

		room.PlayerSet(p)

		log.Println("user added to room:", p.ID)

		for {
			player := room.PlayerGet(p.ID)

			if err = conn.ReadJSON(&event); err != nil {
				log.Println("read:", err)
				continue
			}

			if err := room.PlayerAction(player, event); err != nil {
				log.Println("move player :", err)
				continue
			}
		}
	}))

	log.Fatal(app.Listen(":3001"))
}
