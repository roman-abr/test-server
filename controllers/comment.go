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
)

var cCollection = database.Collection("comment")

type CreateCommentParams struct {
	Text string `json:"text"`
}

func CreateComment(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var params CreateCommentParams
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
	postId, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	comment := models.Comment{
		Id:     primitive.NewObjectID(),
		Text:   params.Text,
		UserId: objID,
		PostId: postId,
	}
	_, err = cCollection.InsertOne(context.Background(), comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err := json.Marshal(comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func GetComments(ctx *fasthttp.RequestCtx) ([]byte, error) {
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"post": objID}
	var comments []models.Comment
	result, err := cCollection.Find(context.Background(), filter)
	result.All(context.Background(), &comments)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err := json.Marshal(comments)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func GetComment(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var comment models.Comment
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	commentID, err := primitive.ObjectIDFromHex(ctx.UserValue("comment_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": commentID, "post": objID}
	err = cCollection.FindOne(context.Background(), filter).Decode(&comment)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Comment not found")
	}
	body, err := json.Marshal(comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}

func RemoveComment(ctx *fasthttp.RequestCtx) ([]byte, error) {
	objID, err := primitive.ObjectIDFromHex(ctx.UserValue("post_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	commentID, err := primitive.ObjectIDFromHex(ctx.UserValue("comment_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(ctx.UserValue("session_user").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": commentID, "post": objID, "user": userId}
	_, err = cCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Post not found")
	}
	return nil, nil
}

func UpdateComment(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var comment models.Comment
	var params CreateCommentParams
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
	commentID, err := primitive.ObjectIDFromHex(ctx.UserValue("comment_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(ctx.UserValue("session_user").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filter := bson.M{"_id": commentID, "post": objID, "user": userId}
	updateQuery := bson.D{{"$set", bson.D{{"text", params.Text}}}}
	err = cCollection.FindOne(context.Background(), filter).Decode(&comment)
	if err != nil {
		log.Println(err)
		return nil, utils.NotFoundError("Comment not found")
	}
	_, err = pCollection.UpdateOne(context.Background(), filter, updateQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	comment.Text = params.Text
	body, err := json.Marshal(comment)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}
