package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Owner     string `json:"owner"`
	OwnerID   string `json:"ownerid"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Title: "Grocery Shopping", Owner: "Steven Weatherspoon", OwnerID: "1", Completed: false},
	{ID: "2", Title: "Programming", Owner: "Steven Weatherspoon", OwnerID: "1", Completed: false},
	{ID: "3", Title: "Test Code", Owner: "Steven Weatherspoon", OwnerID: "1", Completed: false},
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodosByOwner)
	router.POST("/todos/:id", postTodos)

	router.Run()
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func postTodos(c *gin.Context) {
	var newTodo todo

	if err := c.BindJSON(&newTodo); err != nil {
		return
	}

	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodosByOwner(c *gin.Context) {
	id := c.Param("id")
	var response []todo

	for _, todo := range todos {
		if todo.OwnerID == id {
			response = append(response, todo)
		}
	}
	if len(response) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}
