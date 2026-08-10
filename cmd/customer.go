package main

import (
	"log/slog"
	"martabak-tracker-go/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerData struct {
	Ttitle   string
	Order    models.Order
	Statuses []string
}

type OrderFormData struct {
	MartabakTypes []string
	MartabakSizes []string
}

type OrderRequest struct {
	Name          string   `form:"name" binding:"required,min=3,max=50"`
	Phone         string   `form:"phone" binding:"required,numeric,min=8,max=15"`
	Address       string   `form:"address" binding:"required,min=5,max=200"`
	Sizes         []string `form:"sizes" binding:"required,min=1,dive,valid_martabak_size"`
	MartabakTypes []string `form:"martabak" binding:"required,min=1,dive,valid_martabak_type"`
	Instructions  []string `form:"instructions" binding:"max=200"`
}

func (h *Handler) ServeNewOrderForm(c *gin.Context) {
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		MartabakTypes: models.MartabakTypes,
		MartabakSizes: models.MartabakSizes,
	})
}

func (h *Handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderItems := make([]models.OrderItem, len(form.Sizes))
	for i := range orderItems {
		orderItems[i] = models.OrderItem{
			Size:        form.Sizes[i],
			Martabak:    form.MartabakTypes[i],
			Instruction: form.Instructions[i],
		}
	}
	order := models.Order{
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Status:       models.OrderStatuses[0],
		Items:        orderItems,
	}
	if err := h.orders.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "error", err)
		c.String(http.StatusInternalServerError, "Something went wrong while creating the order. Please try again later.")
		return
	}
	slog.Info("Order created successfully", "orderId", order.ID, "customer", order.CustomerName)

	c.Redirect(http.StatusSeeOther, "/customer/"+order.ID)
}

func (h *Handler) serveCustomer(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
	}
	order, err := h.orders.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "Order not found")
		return
	}
	c.HTML(http.StatusOK, "customer.tmpl", CustomerData{
		Ttitle:   "Martabak Enak - Order Details " + order.ID,
		Order:    *order,
		Statuses: models.OrderStatuses,
	})
}
