package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID 			primitive.ObjectID 	`json:"_id,omitempty" bson:"_id,omitempty"`
	Completed 	bool 				`json:"completed"`
	Body 		string 				`json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Started!!!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")

	collection = client.Database("todo_db").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	app.Get("/api/todo", getTodos)
	app.Post("/api/todo", createTodo)
	app.Patch("/api/todo/:id", updateTodo)
	app.Delete("/api/todo/:id", deleteTodo)

	PORT := os.Getenv("PORT")
	if PORT == ""{
		PORT = "5000" 
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	
	if err != nil{
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()){
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if error := c.BodyParser(todo); error != nil {
		return error
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	data, err := collection.InsertOne(context.Background(), todo)
	if err != nil{
		return err
	}

	todo.ID = data.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(id)
	
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "invalid todo ID"})
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": bson.M{"completed": true}})

	if err != nil {
		return c.Status(404).JSON(err)
	}
	
	return c.Status(200).JSON(fiber.Map{"msg": "Successfully updated!!!"})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(id)
	
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "invalid todo ID"})
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objId})

		if err != nil {
		return c.Status(404).JSON(err)
	}

	return c.Status(200).JSON(fiber.Map{"msg": "Successfully deleted!!!"})
}

// func main(){
// 	fmt.Println("Started")

// 	app := fiber.New()

// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	PORT := os.Getenv("PORT")

// 	todos := []Todo{}

// 	// Get all Todos
// 	app.Get("/api/todo", func(c *fiber.Ctx) error {
// 		return c.Status(200).JSON(todos)
// 	})

// 	// Add Todo
// 	app.Post("/api/todo/add", func(c *fiber.Ctx) error {
// 		todo := &Todo{}
// 		if error := c.BodyParser(todo); error != nil {
// 			return error
// 		}

// 		if todo.Body == "" {
// 			return c.Status(400).JSON(fiber.Map{"msg": "Todo body is required"})
// 		}

// 		todo.ID = len(todos) + 1
// 		todos = append(todos, *todo)

// 		return c.Status(200).JSON(todo)
// 	}) 
	
// 	// Update Todo
// 	app.Patch("/api/todo/update/:id", func(c *fiber.Ctx) error {
// 		id := c.Params("id")
// 		for i, todo := range todos{
// 			if fmt.Sprint(todo.ID) == id {
// 				todos[i].Completed = true
// 				return c.Status(200).JSON(todos[i])
// 			}
// 		}
// 		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
// 	})

// 	// Delete Todo
// 	app.Delete("/api/todo/delete/:id", func(c *fiber.Ctx) error {
// 		id := c.Params("id")
		
// 		for i, todo := range todos {
// 			if fmt.Sprint(todo.ID) == id {
// 				todos = slices.Delete(todos, i, i+1)
// 				return c.Status(200).JSON(fiber.Map{"mgs": "Deleted"})
// 			}
// 		}
// 		return c.Status(404).JSON(fiber.Map{"msg": "Todo not found"})
// 	})

// 	log.Fatal(app.Listen(":" + PORT))
// }
