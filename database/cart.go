package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kshzz24/ecomm-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct   = errors.New("cannot find the requested product")
	ErrCantDecodeProduct = errors.New("cannot decode product data")
	ErrUserIdIsNotValid  = errors.New("user id is not valid")
	ErrCantUpdateUser    = errors.New("cannot update user information")
	ErrCantRemoveItem    = errors.New("cannot remove item from cart")
	ErrCantGetItem       = errors.New("cannot get item from database")
	ErrCantBuyCartItem   = errors.New("cannot process purchase from cart")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	err = searchFromDb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{
		Key:   "_id",
		Value: id,
	}}
	update := bson.D{{
		Key: "$push",
		Value: bson.D{
			primitive.E{
				Key: "usercart",
				Value: bson.D{{
					Key:   "$each",
					Value: productCart,
				}},
			},
		},
	}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}
func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil

}
func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.OrderID = primitive.NewObjectID()
	ordercart.OrderedAt = time.Now()
	ordercart.OrderCart = make([]models.ProductUser, 0)
	ordercart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{
		Key:   "path",
		Value: "$usercart",
	}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{
		Key:   "_id",
		Value: "$_id",
	},
		{
			Key:   "total",
			Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}},
		}}}}
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	if err != nil {
		panic(err)

	}

	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}

	var total_price int32

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	ordercart.Price = uint64(total_price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		// return
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	// userCollection.
	if err != nil {
		log.Println(err)
	}

	filter1 := bson.D{
		primitive.E{
			Key:   "_id",
			Value: id,
		},
	}

	update2 := bson.M{
		"$push": bson.M{"orders.$[].order_list": bson.M{
			"$each": getcartitems.UserCart,
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter1, update2)

	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{
		primitive.E{Key: "_id", Value: id},
	}
	update3 := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "usercart",
					Value: usercart_empty,
				},
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)

	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var product_details models.ProductUser

	var orders_detail models.Order

	orders_detail.OrderID = primitive.NewObjectID()
	orders_detail.OrderedAt = time.Now()

	orders_detail.OrderCart = make([]models.ProductUser, 0)

	orders_detail.PaymentMethod.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}

	orders_detail.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
