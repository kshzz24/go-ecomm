package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("cannot find the requested product")
	ErrCantDecodeProduct  = errors.New("cannot decode product data")
	ErrUserIdIsNotValid   = errors.New("user id is not valid")
	ErrCantUpdateUser     = errors.New("cannot update user information")
	ErrCantRemoveItemCart = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from database")
	ErrCantBuyCartItem    = errors.New("cannot process purchase from cart")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
