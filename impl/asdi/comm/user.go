package main

import "github.com/gofiber/fiber/v2"

type User struct {
	Name    string
	Balance int
}

func (app *App) GetUsers(c *fiber.Ctx) error {
	return nil
}

func (app *App) GetUser(c *fiber.Ctx) error {

	name := c.Params("id")
	balance, err := app.db.GetUser(name)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"name":     name,
		"balance:": balance,
	})
}

func (app *App) CreateUser(c *fiber.Ctx) error {

	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	app.db.CreateUser(user.Name, user.Balance)

	return c.JSON(fiber.Map{
		"message": "create user: " + user.Name,
	})
}

func (app *App) DeleteUser(c *fiber.Ctx) error {
	return nil
}
