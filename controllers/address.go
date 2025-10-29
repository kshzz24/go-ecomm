package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kshzz24/ecomm-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addresses models.Address
		addresses.AddressID = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{
			Key: "$match",
			Value: bson.D{bson.E{
				Key:   "_id",
				Value: usert_id,
			}},
		}}
		unwind := bson.D{{
			Key: "$unwind",
			Value: bson.D{
				bson.E{
					Key:   "path",
					Value: "$address",
				},
			},
		}}
		group := bson.D{{
			Key: "$group",
			Value: bson.D{
				bson.E{
					Key:   "_id",
					Value: "$address.address_id",
				},
				bson.E{
					Key: "count",
					Value: bson.D{bson.E{
						Key:   "$sum",
						Value: 1,
					}},
				},
			},
		}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "internal server error")
			return
		}

		var addressInfo []bson.M
		if err = cursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
			if size < 2 {
				filter := bson.D{
					primitive.E{
						Key:   "_id",
						Value: usert_id,
					},
				}
				update := bson.D{
					bson.E{
						Key: "$push",
						Value: bson.D{
							bson.E{
								Key:   "address",
								Value: addresses,
							},
						},
					},
				}

				_, err := UserCollection.UpdateOne(ctx, filter, update)
				if err != nil {
					fmt.Println(err)
				}

			} else {
				c.IndentedJSON(400, "Not Allowed")
			}
		}

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var editaddress models.Address
		if err = c.BindJSON(editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{
			bson.E{
				Key:   "_id",
				Value: usert_id},
		}

		update := bson.D{
			bson.E{
				Key: "$set",
				Value: bson.D{
					bson.E{
						Key:   "address.0.house_name",
						Value: editaddress.House,
					},
					bson.E{
						Key:   "address.0.street_name",
						Value: editaddress.Street,
					},
					bson.E{
						Key:   "address.0.city_name",
						Value: editaddress.City,
					},
					bson.E{
						Key:   "address.0.pincode",
						Value: editaddress.Pincode,
					},
				},
			},
		}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(500, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Updated the home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var editaddress models.Address
		if err = c.BindJSON(editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{
			bson.E{
				Key:   "_id",
				Value: usert_id},
		}
		update := bson.D{
			bson.E{
				Key: "$set",
				Value: bson.D{
					bson.E{
						Key:   "address.1.house_name",
						Value: editaddress.House,
					},
					bson.E{
						Key:   "address.1.street_name",
						Value: editaddress.Street,
					},
					bson.E{
						Key:   "address.1.city_name",
						Value: editaddress.City,
					},
					bson.E{
						Key:   "address.1.pincode",
						Value: editaddress.Pincode,
					},
				},
			},
		}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(500, "something went wrong")
			return
		}

		ctx.Done()
		c.IndentedJSON(200, "Updated the work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{
			bson.E{
				Key:   "_id",
				Value: usert_id},
		}

		update := bson.D{
			{
				Key:   "$set",
				Value: bson.D{bson.E{Key: "address", Value: addresses}},
			},
		}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "wrong command")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, gin.H{
			"message": "successfully Deleted",
		})

	}
}
