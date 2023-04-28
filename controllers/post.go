package controllers

import (
	"context"
	"encoding/json"
	"log"

	"test-server/database"
	"test-server/database/models"
	"test-server/utils"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var pCollection = database.Collection("post")

type CreatePostParams struct {
	Text string `json:"text"`
}

func CreatePost(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var params CreatePostParams
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("session_user").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	post := models.Post{
		Id:     primitive.NewObjectID(),
		Text:   params.Text,
		UserId: objID,
	}
	_, err = pCollection.InsertOne(context.Background(), post)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

type AggregatedPost struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Text       string             `json:"text" bson:"text"`
	User       []models.User      `json:"user"`
	CommsCount int                `json:"comms_count" bson:"comms_count"`
}

func GetPosts(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var aggs []AggregatedPost
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "user"}, {Key: "localField", Value: "user"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "user"}}}},
		{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "comment"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "post_id"}, {Key: "as", Value: "comms"}}}},
		{{Key: "$addFields", Value: bson.D{{Key: "comms_count", Value: bson.D{{Key: "$size", Value: "$comms"}}}}}},
	}
	result, err := pCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result.All(context.Background(), &aggs)
	body, err := json.Marshal(aggs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func GetPost(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var post models.Post
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": objID}
	err = pCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Post not found")
	}
	body, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func RemovePost(ctx *fasthttp.RequestCtx) ([]byte, error) {
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(ctx.UserValue("session_user").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": objID, "user": userId}
	_, err = pCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Post not found")
	}
	return nil, nil
}

func UpdatePost(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var params CreatePostParams
	var post models.Post
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(ctx.UserValue("session_user").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": objID, "user": userId}
	err = pCollection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Post not found")
	}
	updateQuery := bson.D{{"$set", bson.D{{"text", params.Text}}}}
	_, err = pCollection.UpdateOne(context.Background(), filter, updateQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	post.Text = params.Text
	body, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}
