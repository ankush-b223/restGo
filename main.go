package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	qmgo "github.com/qiniu/qmgo"
)

var database *qmgo.Database
var collection *qmgo.Collection

func home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"YES": "ALL GOOD",
	})
}

func main() {

	const dbUri = "mongodb://mongo:27017"
	fmt.Println("Connecting to DB")

	ctx := context.Background()

	connection, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: dbUri})

	database = connection.Database("test")
	collection = database.Collection("books")

	defer func() {
		if err = connection.Close(ctx); err != nil {
			fmt.Println("Closing DB Connection", dbUri)
			panic(err)
		}
	}()

	router := gin.Default()

	router.GET("/", home)
	router.GET("/book", getBooks)
	router.GET("/book/:id", getBook)
	router.PATCH("/book/:id", updateBook)
	router.DELETE("/book/:id", deleteBook)

	router.POST("/book", createBook)

	router.Run(":8080")

}
