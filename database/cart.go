package database

import (
	"context"
	"errors"
	"log"
	"time"
	"github.com/Triptiverma003/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	ErrCantFindProduct = errors.New("can't find product")
	errCantDecodeProducts = errors.New("cant find the product")
	errUserIdIsNotValid = errors.New("not valid user")
	errCantUpdateUser = errors.New("cannot add this product to cart")
	errCantRemoveItemCart = errors.New("cannot remove this items from cart")
	errCantGetItem = errors.New("was unable to get the items from the cart")
	errCantBuyCartItem = errors.New("cannot update the purchase")
	)

func AddProductToCart( ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error{
	searchfromdb, err := prodCollection.Find(ctx , bson.M{"_id": productId})
	if err != nil{
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil{
		log.Println(err)
		return errCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err!= nil{
		log.Println(err)
		return errUserIdIsNotValid
	}

	filter:= bson.D{primitive.E{Key: "_id", Value: id}}
	update:= bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart" , Value: bson.D{{Key: "$each" , Value: productCart}}}}}}

	_, err  = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errCantUpdateUser
	}
	return nil 
}

func RemoveCartItem( ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id , err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.Println(err)
		return errUserIdIsNotValid
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$pull" : bson.M{"usercart" : bson.M{"_id" : productID}}}
	_, err = prodCollection.UpdateMany(ctx , filter , update)
	if err!=nil {
		return errCantRemoveItemCart
	}
	return nil 
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_Id = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true


	unwind := bson.D{{Key: "$unwind" , Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group" , Value: bson.D{primitive.E{Key: "_id", Value: "$_id"} , {Key: "total", Value: bson.D{primitive.E{Key: "$sum" , Value: "$usercart.price"}}}}}}
	currentresults , err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind , grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart);
	err != nil {
		panic(err)
	}
	var total_price int32
	for _ , user_item := range getusercart{
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id" , Value: id}}
	update:= bson.D{{Key: "$push" , Value: bson.D{primitive.E{Key: "orders" , Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil{
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id" , Value: id}}).Decode(&getcartitems)

	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id" , Value : id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}

	_ , err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	usercart_empty := make([]models.ProductUser , 0)
	filter3 := bson.D{primitive.E{Key: "_id" , Value: id}}
	update3 := bson.D{{Key: "$set" , Value: bson.D{primitive.E{Key: "usercart" , Value: usercart_empty}}}}

	_ , err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return errCantBuyCartItem
	}
	return nil
}

func instantBuyer(ctx context.Context , prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error{
	id, err := primitive.ObjectIDFromHex(UserID)
	
	if err != nil{
		log.Println(err)
		return errUserIdIsNotValid
	}
	var product_details models.ProductUser
	var order_details models.Order

	order_details.Order_Id = primitive.NewObjectID()
	order_details.Ordered_At = time.Now()
	order_details.Order_Cart - make([]models.ProductUser , 0)
	orders_details.Payment_Method.COD = true
	err = prodCollection.FindOne(ctx , bson.D{primitive.E{Key: "_id" , Value: productID}}).Decode(&product_details)

	if err != nil {
		log.Panicln(err)
	}
	order_details.Price  = product_details.Price

	filter := 
	update := 
	 _ , err = userCollection.UpdateOne(ctx , filter , update)

	 if err != nil{
		log.Println(err)
	 }
	 filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	 update2 := bson.M{"$push": bson.M{"orders.$[].order_list" : product_details}}

	 _ , err = userCollection.UpdateOne(ctx , filter2 , update2)

	 if err != nil {
		log.Println(err)
	 }
	 return nil
}
