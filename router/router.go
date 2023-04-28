package router

import (
	"log"

	"test-server/controllers"
	"test-server/middleware"

	"github.com/fasthttp/router"
)

func InitRoutes() *router.Router {
	r := router.New()

	// routes
	r.POST("/auth/signin", middleware.Base(controllers.SignIn))
	r.POST("/auth/signup", middleware.Base(controllers.SignUp))

	r.POST("/posts", middleware.Auth(controllers.CreatePost))
	r.GET("/posts", middleware.Auth(controllers.GetPosts))
	r.GET("/posts/{post_id}", middleware.Auth(controllers.GetPost))
	r.DELETE("/posts/{post_id}", middleware.Auth(controllers.RemovePost))
	r.PUT("/posts/{post_id}", middleware.Auth(controllers.UpdatePost))

	r.POST("/posts/{post_id}/comments", middleware.Auth(controllers.CreateComment))
	r.GET("/posts/{post_id}/comments", middleware.Auth(controllers.GetComments))
	r.GET("/posts/{post_id}/comments/{comment_id}", middleware.Auth(controllers.GetComment))
	r.DELETE("/posts/{post_id}/comments/{comment_id}", middleware.Auth(controllers.RemoveComment))
	r.PUT("/posts/{post_id}/comments/{comment_id}", middleware.Auth(controllers.UpdateComment))

	//routes

	for m, paths := range r.List() {
		for _, p := range paths {
			log.Println("Added route: ", m, p)
		}
	}
	return r
}
