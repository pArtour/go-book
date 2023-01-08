package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pArtour/go-book/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// Init validator
var validate = validator.New()

// CreateBook handles POST requests to add new books
func CreateBook(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book model.Book

		// Bind JSON to book struct
		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		if err := validate.Struct(book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set created_at and updated_at
		book.Created_at = time.Now()
		book.Updated_at = time.Now()

		// Set book ID
		book.ID = primitive.NewObjectID()

		// Insert book to database
		result, err := collection.InsertOne(ctx, book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return created book
		c.JSON(http.StatusCreated, gin.H{"book": result})
	}
}

// GetBook handles GET requests to get a book by id
func GetBook(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		bookId := c.Param("book_id")
		var book model.Book

		objectId, _ := primitive.ObjectIDFromHex(bookId)

		err := collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching book."})
			return
		}

		c.JSON(http.StatusOK, book)
	}
}

// UpdateBook handles PATCH requests to update a book by id
func UpdateBook(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get book ID from request
		id := c.Param("id")
		var book model.Book

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert book ID to primitive.ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"_id": objectID}
		var updateObj primitive.D

		if book.Title != nil {
			updateObj = append(updateObj, bson.E{"title", book.Title})
		}

		if book.Author != nil {
			updateObj = append(updateObj, bson.E{"author", book.Author})
		}

		if book.Description != nil {
			updateObj = append(updateObj, bson.E{"description", book.Description})
		}

		book.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: book.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err = collection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"book": book})
	}
}

// DeleteBook handles DELETE requests to delete a book by id
func DeleteBook(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get book ID from request
		id := c.Param("id")

		// Convert book ID to primitive.ObjectID
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Delete book from database
		_, err = collection.DeleteOne(ctx, model.Book{ID: objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return deleted book
		c.JSON(http.StatusOK, gin.H{"message": "Book item deleted successfully."})
	}
}

// GetBooks handles GET requests to get all books
func GetBooks(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching book list"})
			return
		}

		var allBooks []bson.M
		if err := result.All(ctx, &allBooks); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allBooks)
	}
}
