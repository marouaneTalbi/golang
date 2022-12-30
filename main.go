package main

import (
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
	r.POST("/products", controllers.ProductsCreate)
	r.GET("/products", controllers.ProductsIndex)
	r.GET("/products/:id", controllers.ProductsShow)
	r.PUT("/products/:id", controllers.ProductsUpdate)
	r.DELETE("/products/:id", controllers.ProductDelete)

	r.POST("/payments", controllers.PaymentCreate)
	r.GET("/payments", controllers.PaymentIndex)
	r.GET("/payments/:id", controllers.PaymentShow)
	r.PUT("/payments/:id", controllers.PaymentUpdate)
	r.DELETE("/payments/:id", controllers.PaymentDelete)

	r.Run()
}
