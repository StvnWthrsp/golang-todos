package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Todo struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Task    string             `json:"task,omitempty"`
	Owner   string             `json:"owner,omitempty"`
	OwnerID float64            `json:"ownerid,omitempty"`
	Status  bool               `json:"status,omitempty"`
}

const uri = "mongodb://root:rootpw@localhost:27017/?authSource=admin"

var todoCollection *mongo.Collection
var validate = validator.New()

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	todoCollection = client.Database("Todo").Collection("todos")

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.POST("/todos", postTodos)

	router.Run()

}

func getTodos(c *gin.Context) {
	var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var todos []bson.M

	cursor, err := todoCollection.Find(context, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = cursor.All(context, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cancel()

	fmt.Println(todos)

	c.JSON(http.StatusOK, todos)
}

func postTodos(c *gin.Context) {
	var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var newTodo Todo

	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(newTodo)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	newTodo.ID = primitive.NewObjectID()

	result, insertErr := todoCollection.InsertOne(context, newTodo)
	if insertErr != nil {
		msg := fmt.Sprintf("Todo was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result)
}
