package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID 			int 	`json: id`
	Completed 	bool 	`json: completed`
	Body 		string 	`json: body`
}

func main(){
	fmt.Println("Started")

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	// Get all Todos
	app.Get("/api/todo", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Add Todo
	app.Post("/api/todo/add", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if error := c.BodyParser(todo); error != nil {
			return error
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(200).JSON(todo)
	}) 
	
	// Update Todo
	app.Patch("/api/todo/update/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos{
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})

	// Delete Todo
	app.Delete("/api/todo/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = slices.Delete(todos, i, i+1)
				return c.Status(200).JSON(fiber.Map{"mgs": "Deleted"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
