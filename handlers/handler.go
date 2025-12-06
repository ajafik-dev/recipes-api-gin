package handlers

import (
	"fmt"
	"net/http"
	"recipes-api-gin/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.
	Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.
	Context) {

	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting a new recipe"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": recipe.Name, "tags": recipe.Tags, "ingredients": recipe.Ingredients, "instructions": recipe.Instructions}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "Recipe has been Updated"})

}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
//   - application/json
//
// responses:
// '200':
//
//	description: Successful operation
func (handler *RecipesHandler) ListRecipesHandler(c *gin.
	Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)
	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
}

func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cur, err := handler.collection.Find(handler.ctx, bson.M{"tags": bson.M{"$regex": tag, "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	listOfRecipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		listOfRecipes = append(listOfRecipes, recipe)
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
