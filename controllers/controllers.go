package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/kshzz24/ecomm-go/database"
	"github.com/kshzz24/ecomm-go/models"
	generate "github.com/kshzz24/ecomm-go/tokens"

	// token "github.com/kshzz24/ecomm-go/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()
var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or password is incorrect"
		valid = false
	}

	return valid, msg

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		validationErr := Validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(user.Password, founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(founduser.Email, founduser.FirstName, founduser.LastName, founduser.UserID)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.UserID)
		c.JSON(http.StatusFound, founduser)

	}
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := Validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exists",
			})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{
			"phone": user.Phone,
		})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "this phone no. is already in use",
			})
			return
		}

		password := HashPassword(*&user.Password)
		user.Password = password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		token, refreshtoken, _ := generate.TokenGenerator(user.Email, user.FirstName, user.LastName, user.UserID)

		user.Token = token
		user.RefreshToken = refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)

		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "user not get created",
			})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed in")

	}
}

// func ProductViewerAdmin() gin.HandlerFunc {

// }

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		err = cursor.All(ctx, &productlist)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		c.IndentedJSON(200, productlist)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

		var searchproducts []models.Product

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("Query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invalid search index",
			})

			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQuerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			c.IndentedJSON(404, "something went wrong while fetching data")
			return
		}

		err = searchQuerydb.All(ctx, &searchproducts)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchQuerydb.Close(ctx)

		if err := searchQuerydb; err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid req")
			return
		}

		c.IndentedJSON(200, searchproducts)

	}
}

//

//

//
