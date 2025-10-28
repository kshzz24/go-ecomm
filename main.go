package main

import (
	"log"
	"os"

	"github.com/kshzz24/ecomm-go/controllers"
	"github.com/kshzz24/ecomm-go/database"
	"github.com/kshzz24/ecomm-go/middleware"
	"github.com/kshzz24/ecomm-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())


	router.GET("/addtocard", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/chartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.Tnstantbuy())

	log.Fatal(router.Run(":"+port))




}
