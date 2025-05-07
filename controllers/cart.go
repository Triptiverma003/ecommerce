package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	database "github.com/Triptiverma003/ecommerce-go/DataBase"
	"github.com/Triptiverma003/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}
func NewApplication(prodCollection, userCollection *mongo.Collection)*Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddtoCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")

			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID , err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		var ctx , cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.UserCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200 , "Successfully added to cart")
	}

}

func (app *Application)RemoveItem() gin.HandlerFunc{
	return func(c*gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")

			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		var ctx , cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
			return 
		}
		c.IndentedJSON(200 , "Successfully removed from cart")

	}

}

func GetItemFromCart() gin.HandlerFunc{
	return func (c *gin.Context){
		c.Query("id")
		user_id:=c.Query("id") 

		if user_id == ""{
			c.Header("Content-Type" , "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error" : "invalid id"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		
		var ctx , cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		 var filledcart models.User
		 err := UserCollection.FindOne(ctx , bson.D{primitive.E{Key: "_id" , Value: usert_id}}).Decode(&filledcart)

		 if err!= nil{
			log.Println(err)
			c.IndentedJSON(500 , "not found")
			return 
		 }

		 filter_match := bson.D{{Key: "$match" , Value: bson.D{primitive.E{Key: "_id" , Value: usert_id}}}}
		 unwind := bson.D{{Key: "$unwind" , Value: bson.D{primitive.E{Key: "path" , Value: "$usercart"}}}}
		 grouping := bson.D{{Key : "$group" , Value:bson.D{primitive.E{Key: "_id" , Value: "$_id"}, {Key: "total" , Value:bson.D{primitive.E{Key: "$sum" , Value: "$usercart.price"}}}}}}
		 pointCursor , err := UserCollection.Aggregate(ctx , mongo.Pipeline{filter_match , unwind , grouping})

		 if err!=nil{
			log.Println(err)
		 }
		 var listing []bson.M
		 if err:= pointCursor.All(ctx , &listing ); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		 }

		 for _,json := range listing{
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200 , filledcart.UserCart)
		 }
		 ctx.Done()
	}
}

func (app *Application)BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID:= c.Query("id")
		if userQueryID == "" {
			log.Panicln("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err:= database.BuyItemFromCart(ctx, app.prodCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "successfully purchased the item")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product id is empty")

			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		var ctx , cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		database.instantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200 , "Successfully purchased the item")
	}	
}