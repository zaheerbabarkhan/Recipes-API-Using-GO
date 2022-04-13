package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zaheerbabarkhan/Recipes-API-Using-GO/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var err error
var client *mongo.Client
var ctx context.Context
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.Run()
}

func init() {

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017/recipesdb?authSource=admin"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("connected Successfully")
	collection = client.Database("recipesdb").Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)

}
