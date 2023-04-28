package main

import (
	"log"

	"test-server/router"

	"github.com/valyala/fasthttp"
)

func main() {
	port := ":3000"
	r := router.InitRoutes()
	log.Printf("Server started on %s", port)
	log.Fatalln(fasthttp.ListenAndServe(port, r.Handler))
}
