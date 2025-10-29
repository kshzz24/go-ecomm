package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kshzz24/ecomm-go/database"
	"github.com/kshzz24/ecomm-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully added to card")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully removed item from card")

	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "invalid id",
			})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User

		err := UserCollection.FindOne(ctx, bson.D{{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		filterMatch := bson.D{{
			Key:   "$match",
			Value: bson.D{{Key: "_id", Value: usert_id}},
		}}

		unwind := bson.D{{
			Key:   "$unwind",
			Value: bson.D{{Key: "path", Value: "$usercart"}},
		}}

		grouping := bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}},
			},
		}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unwind, grouping})
		if err != nil {
			log.Println(err)
			return
		}
		var listing []bson.M
		if err = cursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}

		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Panic("User id empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "Successfully placed the order")

	}
}

func (app *Application) Instantbuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User id empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully placed the order")
	}
}
