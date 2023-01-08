package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pArtour/go-book/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookRoutes is the route for books
func BookRoutes(routes *gin.Engine, collection *mongo.Collection) {
	routes.GET("/books", controller.GetBooks(collection))
	routes.POST("/books/create", controller.CreateBook(collection))
	routes.GET("/books/:id", controller.GetBook(collection))
	routes.PATCH("/books/:id", controller.UpdateBook(collection))
	routes.DELETE("/books/:id", controller.DeleteBook(collection))

}
