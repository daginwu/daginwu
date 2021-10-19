package main

import "github.com/gofiber/fiber/v2"

type Txn struct {
	From    User `json:"from"`
	To      User `to:"from"`
	Payload int  `payload:"from"`
}

func (app *App) CreateTxn(c *fiber.Ctx) error {
	return nil
}
