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
	"golang.org/x/crypto/bcrypt"
)

var uCollection = database.Collection("user")

type SignInParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpParams struct {
	SignInParams
	Name string `json:"name"`
}

func SignIn(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var params SignInParams
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var user models.User
	filter := bson.D{{"username", params.Username}}
	err = uCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, utils.BadRequestError("Invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		log.Println(err)
		return nil, utils.BadRequestError("Invalid username or password")
	}
	token := utils.GenerateToken(user.Id)
	body, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func SignUp(ctx *fasthttp.RequestCtx) ([]byte, error) {
	var params SignUpParams
	err := json.Unmarshal(ctx.PostBody(), &params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     params.Name,
		Username: params.Username,
		Password: string(password),
	}
	result, err := uCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newUser.Id = result.InsertedID.(primitive.ObjectID)
	body, err := json.Marshal(newUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return body, nil
}
