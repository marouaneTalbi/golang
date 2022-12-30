package controllers

import (
	"example/golang_/initializers"
	"example/golang_/models"

	"github.com/gin-gonic/gin"
)

func ProductsCreate(c *gin.Context) {

	var newProduct struct {
		Name  string
		Price float64
	}

	c.Bind(&newProduct)

	product := models.Product{Name: newProduct.Name, Price: newProduct.Price}
	result := initializers.DB.Create(&product)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create product",
			"error":   newProduct,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": product,
	})
}

func ProductsIndex(c *gin.Context) {
	var products []models.Product
	initializers.DB.Find(&products)

	c.JSON(200, gin.H{
		"message": products,
	})
}

func ProductsShow(c *gin.Context) {
	var product models.Product
	initializers.DB.First(&product, c.Param("id"))

	c.JSON(200, gin.H{
		"product": product,
	})
}

func ProductsUpdate(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct struct {
		Name  string
		Price float64
	}
	c.Bind(&updatedProduct)
	var product models.Product
	initializers.DB.First(&product, id)

	initializers.DB.Model(&product).Updates(models.Product{Name: updatedProduct.Name, Price: updatedProduct.Price})

	c.JSON(200, gin.H{
		"product": product,
	})
}

func ProductDelete(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	initializers.DB.First(&product, id)

	initializers.DB.Delete(&product)

	c.JSON(200, gin.H{
		"message": "Product deleted",
	})
}
