package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName      string             `bson:"first_name,omitempty" json:"first_name,omitempty" validate:"required,min=2,max=30"`
	LastName       string             `bson:"last_name,omitempty" json:"last_name,omitempty" validate:"required,min=2,max=30"`
	Password       string             `bson:"password,omitempty" json:"password,omitempty" validate:"required,min=8"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Phone          string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Token          string             `bson:"token,omitempty" json:"token,omitempty"`
	RefreshToken   string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	UserID         string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserCart       []ProductUser      `bson:"usercart,omitempty" json:"usercart,omitempty"`
	AddressDetails []Address          `bson:"address,omitempty" json:"address,omitempty"`
	OrderStatus    []Order            `bson:"orders,omitempty" json:"orders,omitempty"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductName string             `bson:"product_name,omitempty" json:"product_name,omitempty"`
	Price       uint64             `bson:"price,omitempty" json:"price,omitempty"`
	Rating      uint               `bson:"rating,omitempty" json:"rating,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"product_id,omitempty" json:"product_id,omitempty"`
	ProductName string             `bson:"product_name,omitempty" json:"product_name,omitempty"`
	Price       uint64             `bson:"price,omitempty" json:"price,omitempty"`
	Rating      uint               `bson:"rating,omitempty" json:"rating,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	House     string             `bson:"house_name,omitempty" json:"house_name,omitempty"`
	Street    string             `bson:"street_name,omitempty" json:"street_name,omitempty"`
	City      string             `bson:"city_name,omitempty" json:"city_name,omitempty"`
	Pincode   uint16             `bson:"pin_code,omitempty" json:"pin_code,omitempty"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderCart     []ProductUser      `bson:"order_list,omitempty" json:"order_list,omitempty"`
	OrderedAt     time.Time          `bson:"ordered_at,omitempty" json:"ordered_at,omitempty"`
	Price         uint64             `bson:"total_price,omitempty" json:"total_price,omitempty"`
	Discount      uint64             `bson:"discount,omitempty" json:"discount,omitempty"`
	PaymentMethod Payment            `bson:"payment_method,omitempty" json:"payment_method,omitempty"`
}

type Payment struct {
	Digital bool `bson:"digital,omitempty" json:"digital,omitempty"`
	COD     bool `bson:"cod,omitempty" json:"cod,omitempty"`
}
