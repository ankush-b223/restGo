package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	field.DefaultField `bson:",inline"`
	Title              string `bson:"title" validate:"required"`
	Author             string `bson:"author"`
}

type BookCreateUpdateRequest struct {
	Title  string `form:"title" binding:"required"`
	Author string `form:"author"`
}

type BookResponse struct {
	Id        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	Author    string             `json:"author"`
	CreatedAt time.Time          `json:"createdAt"  binding:"required"`
	UpdatedAt time.Time          `json:"updatedAt"  binding:"required"`
}

type BookList struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Title string             `json:"title"`
}

func createBook(c *gin.Context) {

	var newBook BookCreateUpdateRequest

	if err := c.ShouldBind(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}

	book := Book{
		Title:  newBook.Title,
		Author: newBook.Author,
	}

	_, err := collection.InsertOne(c, &book)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Retry please")
		return
	}

	c.IndentedJSON(http.StatusCreated, getBooksResponse(book))

}

func getBooksResponse(book Book) (res BookResponse) {

	res = BookResponse{
		book.Id,
		book.Title,
		book.Author,
		book.CreateAt,
		book.UpdateAt,
	}

	return

}

func getBooks(c *gin.Context) {

	var books []BookList

	err := collection.Find(c, bson.M{}).All(&books)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "try again")
		return
	}

	for _, book := range books {
		fmt.Printf("Book ID: %v, Title: %s\n", book.Id, book.Title)
	}

	c.IndentedJSON(http.StatusOK, books)

}

func getBook(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "No bookId found")
		return
	}

	var book Book

	err = collection.Find(c, bson.M{"_id": id}).One(&book)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
		return
	}

	c.IndentedJSON(http.StatusFound, getBooksResponse(book))

}

func deleteBook(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "No bookId found")
		return
	}

	var book Book
	err = collection.Find(c, bson.M{"_id": id}).One(&book)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
		return
	}

	err = collection.RemoveId(c, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
		return
	}

	c.IndentedJSON(http.StatusOK, true)

}

func updateBook(c *gin.Context) {

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "No bookId found")
		return
	}

	var newBook BookCreateUpdateRequest //new book with only to be changed things

	if err := c.ShouldBind(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, "no binding")
		return
	}

	var book Book //the oldBook
	err = collection.Find(c, bson.M{"_id": id}).One(&book)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "old book nt found")
		return
	}

	book.Author = newBook.Author
	book.Title = newBook.Title

	// update in database
	err = collection.ReplaceOne(c, bson.M{"_id": id}, &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something went wrong, Try again after sometime")
		return
	}

	// to send success response on completion
	c.JSON(http.StatusOK, getBooksResponse(book))

}
