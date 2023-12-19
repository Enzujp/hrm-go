package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/options"
)

type MongoInstance struct {
	Client
	Db
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb+srv://jay:jessedavid@menu.08vbsfz.mongodb.net/"

type Employee struct {
	ID string
	Name string
	Salary float64
	Age float64
}

func Connect() error{ // returns error if it goes wrong
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	db := client.Database(dbName)
}

func main() {
	app := fiber.New()

	if err:= Connect(); err != nil {
		log.Fatal("Error trying to connect")
	}

	app.Get("/employee", func(c *fiber.Ctx) error{

	})
	app.Post("/employee")
	app.Put("/employee/:id")
	app.Delete("/employee/:id")
}