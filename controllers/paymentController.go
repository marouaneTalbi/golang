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
	UserId string
	Text   string
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
			c.SSEvent("message", " "+serviceMsg.UserId+" → "+serviceMsg.Text)
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
	result := initializers.DB.Create(&payment)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Failed to pay for product",
			"error":   newPayment,
		})
		return
	}

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Payement is created",
	})

	c.JSON(200, gin.H{
		"message": payment,
	})
}

func (adapter *ginAdapter) PaymentIndex(c *gin.Context) {
	var payments []models.Payment
	initializers.DB.Find(&payments)

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Payement price is ??? ",
	})

	c.JSON(200, gin.H{
		"message": payments,
	})
}

func (adapter *ginAdapter) PaymentShow(c *gin.Context) {
	var payment models.Payment
	initializers.DB.First(&payment, c.Param("id"))
	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Payements are fetched",
	})

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

	c.JSON(200, gin.H{
		"product": payment,
	})

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text:   "Payement is updated",
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
		UserId: "1",
		Text:   "Payement is deleted",
	})
}
