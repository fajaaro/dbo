package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fajaaro/dbo/app"
	"github.com/fajaaro/dbo/app/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

type ReqOrder struct {
	CustomerID    uint    `binding:"required" json:"customer_id"`
	ProductName   string  `binding:"required" json:"product_name"`
	Quantity      int     `binding:"required,min=1" json:"quantity"`
	TotalPrice    float64 `binding:"required,min=0" json:"total_price"`
	PaymentStatus string  `binding:"required" json:"payment_status"`
}

func handleValidationError(err error, c *gin.Context) {
	validationErrors := err.(validator.ValidationErrors)
	errorMsg := validationErrors[0].Field() + " not valid"
	res := models.JsonResponse{
		Success: false,
		Error:   &errorMsg,
	}
	c.JSON(http.StatusBadRequest, res)
	c.Abort()
}

func validateOrderInputAndCustomerExistence(req *ReqOrder, db *gorm.DB, res *models.JsonResponse, c *gin.Context) error {
	var customer models.Customer
	result := db.First(&customer, req.CustomerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Customer not found"
			return errors.New(errorMsg)
		}
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return errors.New(errorMsg)
	}

	req.PaymentStatus = strings.ToLower(req.PaymentStatus)
	if req.PaymentStatus != "paid" && req.PaymentStatus != "unpaid" {
		errorMsg := "Invalid payment status"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		return errors.New(errorMsg)
	}

	return nil
}

func OrderController() *OrderRepo {
	return &OrderRepo{DB: app.InitDb()}
}

func (repo *OrderRepo) GetAllOrders(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	limit := c.DefaultQuery("limit", "10")
	limitNum, _ := strconv.Atoi(limit)
	search := c.DefaultQuery("search", "")
	search = strings.ToLower(search)

	var orders []models.Order
	query := repo.DB.Model(&models.Order{})

	if search != "" {
		query = query.Where("product_name ILIKE ?", "%"+search+"%")
	}

	var count int64
	query.Count(&count)

	query = query.Offset((pageNum - 1) * limitNum).Limit(limitNum).Find(&orders)

	if query.Error != nil {
		errorMsg := query.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = map[string]interface{}{
		"orders": orders,
		"count":  count,
	}
	c.JSON(http.StatusOK, res)
}

func (repo *OrderRepo) GetOrderDetail(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	orderID := c.Param("id")

	var order models.Order
	result := repo.DB.First(&order, orderID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Order not found"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusNotFound, res)
			return
		}
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = order
	c.JSON(http.StatusOK, res)
}

func (repo *OrderRepo) InsertOrder(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := ReqOrder{}
	err := c.BindJSON(&req)
	if err != nil {
		handleValidationError(err, c)
		return
	}

	if err := validateOrderInputAndCustomerExistence(&req, repo.DB, &res, c); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	order := &models.Order{
		CustomerID:    req.CustomerID,
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		TotalPrice:    req.TotalPrice,
		PaymentStatus: req.PaymentStatus,
	}
	if req.PaymentStatus == "paid" {
		paidAt := time.Now()
		order.PaidAt = &paidAt
	}
	result := repo.DB.Create(&order)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = order
	c.JSON(http.StatusCreated, res)
}

func (repo *OrderRepo) UpdateOrder(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	orderID := c.Param("id")

	var order models.Order
	result := repo.DB.First(&order, orderID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Order not found"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusNotFound, res)
			return
		}
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	req := ReqOrder{}
	err := c.BindJSON(&req)
	if err != nil {
		handleValidationError(err, c)
		return
	}

	if err := validateOrderInputAndCustomerExistence(&req, repo.DB, &res, c); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	order.ProductName = req.ProductName
	order.Quantity = req.Quantity
	order.TotalPrice = req.TotalPrice
	order.PaymentStatus = req.PaymentStatus
	if req.PaymentStatus == "paid" {
		paidAt := time.Now()
		order.PaidAt = &paidAt
	} else {
		order.PaidAt = nil
	}

	result = repo.DB.Save(&order)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = order
	c.JSON(http.StatusOK, res)
}

func (repo *OrderRepo) DeleteOrder(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	orderID := c.Param("id")

	var order models.Order
	result := repo.DB.First(&order, orderID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Order not found"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusNotFound, res)
			return
		}
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	result = repo.DB.Delete(&order)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = "Order deleted successfully"
	c.JSON(http.StatusOK, res)
}
