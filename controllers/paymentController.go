package controllers

import (
	"example/golang_/initializers"
	"example/golang_/models"

	"github.com/gin-gonic/gin"
)

func PaymentCreate(c *gin.Context) {

	var newPayment struct {
		ProductID int64
		PricePaid float64
	}

	c.Bind(&newPayment)

	payment := models.Payment{ProductID: newPayment.ProductID, PricePaid: newPayment.PricePaid}
	result := initializers.DB.Create(&payment)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Failed to pay for product",
			"error":   newPayment,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": payment,
	})
}

func PaymentIndex(c *gin.Context) {
	var payments []models.Payment
	initializers.DB.Find(&payments)

	c.JSON(200, gin.H{
		"message": payments,
	})
}

func PaymentShow(c *gin.Context) {
	var payment models.Payment
	initializers.DB.First(&payment, c.Param("id"))

	c.JSON(200, gin.H{
		"product": payment,
	})
}

func PaymentUpdate(c *gin.Context) {
	id := c.Param("id")
	var updatedPayment struct {
		ProductID int64
		PricePaid float64
	}
	c.Bind(&updatedPayment)
	var payment models.Payment
	initializers.DB.First(&payment, id)

	initializers.DB.Model(&payment).Updates(models.Payment{ProductID: updatedPayment.ProductID, PricePaid: updatedPayment.PricePaid})

	c.JSON(200, gin.H{
		"product": payment,
	})
}

func PaymentDelete(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment
	initializers.DB.First(&payment, id)

	initializers.DB.Delete(&payment)

	c.JSON(200, gin.H{
		"message": "Payment deleted",
	})
}
