package main

import (
	"martabak-tracker-go/internal/models"
)

type Handler struct {
	orders *models.OrderModel
	users  *models.UserModel
}

func NewHandler(dbModel *models.DBModel) *Handler {
	return &Handler{
		orders: &dbModel.Order,
		users:  &dbModel.User,
	}
}
