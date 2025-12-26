package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	EmailId   string             `bson:"emailId" json:"emailId"`
	Password  string             `bson:"password" json:"-"`
	Salt      string             `bson:"salt" json:"-"`
}
