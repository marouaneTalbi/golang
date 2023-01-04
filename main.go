package main

import (
	"example/golang_/broadcasters"
	"example/golang_/controllers"
	"example/golang_/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	b := broadcasters.NewBroadcaster(20)

	ginAdapter := controllers.NewGinAdapter(b)

	r.POST("/products", ginAdapter.ProductsCreate)
	r.GET("/products", ginAdapter.ProductsIndex)
	r.GET("/products/:id", ginAdapter.ProductsShow)
	r.PUT("/products/:id", ginAdapter.ProductsUpdate)
	r.DELETE("/products/:id", ginAdapter.ProductDelete)

	r.POST("/payments", ginAdapter.PaymentCreate)
	r.GET("/payments", ginAdapter.PaymentIndex)
	r.GET("/payments/:id", ginAdapter.PaymentShow)
	r.PUT("/payments/:id", ginAdapter.PaymentUpdate)
	r.DELETE("/payments/:id", ginAdapter.PaymentDelete)

	r.GET("/stream", ginAdapter.Stream)

	r.Run()
}
