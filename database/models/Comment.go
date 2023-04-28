package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Text   string             `json:"text"`
	UserId primitive.ObjectID `json:"user_id" bson:"user"`
	PostId primitive.ObjectID `json:"post_id" bson:"post"`
}
