package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Text   string             `json:"text" bson:"text"`
	UserId primitive.ObjectID `json:"user_id" bson:"user"`
}
