package database

import (
	"errors"
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

func AddProductToCart() {
	// Add product to cart logic here
}

func RemoveCartItem() error {
	// Example usage of errCantFindProduct
}

func BuyItemFromCart(){

}

func instantBuyer(){

}
