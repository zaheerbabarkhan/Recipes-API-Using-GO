package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaheerbabarkhan/Recipes-API-Using-GO/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
	err        error
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}
func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]models.Recipe, 0)
	result, err := handler.collection.Find(handler.ctx, bson.D{
		{Key: "tags", Value: bson.D{
			{Key: "$all", Value: bson.A{tag}},
		}},
	})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	for result.Next(handler.ctx) {
		var recipe models.Recipe
		result.Decode(&recipe)
		listOfRecipes = append(listOfRecipes, recipe)
	}
	if len(listOfRecipes) != 0 {
		c.JSON(http.StatusOK, listOfRecipes)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "No Recipe Found",
	})
}
func (handler *RecipesHandler) GetRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var recipe models.Recipe
	result.Decode(&recipe)
	if len(recipe.ID) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Recipe Found",
		})
		return
	}
	c.JSON(http.StatusOK, recipe)
}
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result.DeletedCount,
	})
}
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if error := c.ShouldBindJSON(&recipe); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, handler.err = handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: recipe.Name},
		{Key: "instructions", Value: recipe.Instructions},
		{Key: "ingredients", Value: recipe.Ingredients},
		{Key: "tags", Value: recipe.Tags},
	}}})
	if handler.err != nil {
		log.Fatal(handler.err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": handler.err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated",
	})
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if error := c.ShouldBindJSON(&recipe); error != nil {
		fmt.Println("ERROR ", error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, handler.err = handler.collection.InsertOne(handler.ctx, recipe)
	if handler.err != nil {
		log.Fatal(handler.err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error While Inserting New Recipe",
		})
	}
	c.JSON(http.StatusOK, recipe)
}
