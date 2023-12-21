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
	Client *mongo.CLient
	Db *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb+srv://jay:jessedavid@menu.08vbsfz.mongodb.net/"

type Employee struct {
	ID string `json: "id,omitempty" bson:"_id,omitempty"` // golang understands json as id, bson for mongo as _id
	Name string `json: "name"`
	Salary float64 `json:"salary" `
	Age float64 `json: "age"`
}

func Connect() error{ // returns error if it goes wrong
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = mongo.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return  err
	}

	mg = MongoInstance{
		Client: client,
		Db: db,
	}
	return nil
}

func main() {
	app := fiber.New()

	if err:= Connect(); err != nil {
		log.Fatal("Error trying to connect")
	}

	app.Get("/employee", func(c *fiber.Ctx) error{ // get all employees
		query := bson.D{{}} // sending empty brackets so as to fetch all employees from db
		cursor, err:= mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employees []Employee = make([]Employee, 0) // using a slice so it returns the employees as objects containing their defining attributes
		if err:= cursor.All(c.Context(), &employees); err != nil {
			return  c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees)
	})


	app.Post("/employee", func(c *fiber.Ctx) error { // using fiber grants us access to our response and request
		collection := mg.Db.Collection("employees")

		employee := new(Employee)

		if err:= c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		} // parse json body

		employee.ID = ""
		insertionResult, err :=collection.InsertOne(c.Context(), employee); if err != nil {// insert one adds one data
			return c.Status(500).SendString(err.Error())}

		filter := bson.D{{key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee) // because golang doesnt understand json
		return c.Status(201).JSON(createdEmployee)
		
	})


	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		employeeID, err := primitive.ObjectIdFromHex(idParam); if err != nil {
			return c.Status(400).SendString(err.Error())}

		employee := new(Employee)

		if err:= c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error)
		}

		query := bson.D{{Key:"_id", Value: employeeID}}
		update:= bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key:"name", Value: employee.Name},
					{Key:"age", Value: employee.Age},
					{Key:"salary", Value: employee.Salary},
				},
			},
		}

		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()

		if err != nil {
			// means filter didnt match any docs
			if err == mongo.ErrNoDocuments {
				 return c.SendStatus(400)
			}
			
			return c.SendStatus(500)
		}

		employeee.ID = idParam
		return c.sendStatus(200).JSON(employee)

	})


	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		employeeId, err := primitive.ObjectIdFromHex(c.Params("id"))
		if err != nil {
			c.SendStatus(400) // probably an invalid id sent 
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), &query)
		
		if err != nil {
			return c.SendStatus(500)
		}

		if result.DeletedCount < 1 { // deleted count is a mongodb function
			c.SendStatus(404)
		}

		return c.status(200).JSON("record deleted")
	})



Log.Fatal(app.Listen(":3000"))