package controllers

import (
	"example/golang_/broadcasters"
	"example/golang_/initializers"
	"example/golang_/models"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type ginAdapter struct {
	broadcaster broadcasters.Broadcaster
}

type Message struct {
	Text string
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewGinAdapter(broadcaster broadcasters.Broadcaster) *ginAdapter {
	return &ginAdapter{
		broadcaster: broadcaster,
	}
}

func (adapter *ginAdapter) Stream(c *gin.Context) {

	listener := make(chan interface{})

	adapter.broadcaster.Register(listener)
	defer adapter.broadcaster.Unregister(listener)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				fmt.Println("not a message")
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", serviceMsg.Text)
			return true
		}
	})

	fmt.Println("stream is OK")
}

func (adapter *ginAdapter) PaymentCreate(c *gin.Context) {

	var newPayment struct {
		ProductID int64
		PricePaid float64
	}

	c.Bind(&newPayment)

	payment := models.Payment{ProductID: newPayment.ProductID, PricePaid: newPayment.PricePaid}

	var product models.Product
	initializers.DB.First(&product, newPayment.ProductID)

	if product.Price == newPayment.PricePaid {
		result := initializers.DB.Create(&payment)
		if result.Error != nil {
			c.JSON(500, gin.H{
				"message": "Failed to pay for product",
				"error":   newPayment,
			})
			return
		}
	} else {
		c.JSON(500, gin.H{
			"message": "Unvalid Price",
			"error":   newPayment,
		})
		return
	}

	b := adapter.broadcaster

	b.Submit(Message{
		Text: fmt.Sprintf("prix: %v", newPayment.PricePaid) + fmt.Sprintf(", Product name: %v", product.Name),
	})

	// c.JSON(200, gin.H{
	// 	"message": payment,
	// })
}

func (adapter *ginAdapter) PaymentIndex(c *gin.Context) {
	var payments []models.Payment
	initializers.DB.Find(&payments)

	// b := adapter.broadcaster
	// fmt.Println(payments)

	// b.Submit(Message{
	// 	Text: "All Payements are fitched",
	// })

	c.JSON(200, gin.H{
		"message": payments,
	})
}

func (adapter *ginAdapter) PaymentShow(c *gin.Context) {
	var payment models.Payment
	initializers.DB.First(&payment, c.Param("id"))
	// b := adapter.broadcaster

	// b.Submit(Message{
	// 	Text: "Payement are fetched : " + fmt.Sprintf("prix: %v", payment.PricePaid),
	// })

	c.JSON(200, gin.H{
		"product": payment,
	})

}

func (adapter *ginAdapter) PaymentUpdate(c *gin.Context) {
	id := c.Param("id")
	var updatedPayment struct {
		ProductID int64
		PricePaid float64
	}

	c.Bind(&updatedPayment)
	var payment models.Payment
	initializers.DB.First(&payment, id)
	initializers.DB.Model(&payment).Updates(models.Payment{ProductID: updatedPayment.ProductID, PricePaid: updatedPayment.PricePaid})

	var product models.Product
	initializers.DB.First(&product, updatedPayment.ProductID)

	c.JSON(200, gin.H{
		"product": payment,
	})

	b := adapter.broadcaster

	b.Submit(Message{
		Text: "Payement is updated" + fmt.Sprintf(", Price: %v", payment.PricePaid) + fmt.Sprintf(", Product name: : %v", product.Name),
	})
}

func (adapter *ginAdapter) PaymentDelete(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment
	initializers.DB.First(&payment, id)

	initializers.DB.Delete(&payment)

	c.JSON(200, gin.H{
		"message": "Payment deleted",
	})

	b := adapter.broadcaster

	b.Submit(Message{
		Text: "Payement is deleted",
	})
}
