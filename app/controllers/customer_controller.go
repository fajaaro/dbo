package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fajaaro/dbo/app"
	"github.com/fajaaro/dbo/app/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	DB *gorm.DB
}

type ReqCustomer struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=7,max=14"`
	Gender      string `json:"gender" binding:"required"`
}

type ReqUpdateCustomer struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,min=7,max=14"`
	Gender      string `json:"gender" binding:"required"`
}

func CustomerController() *CustomerRepo {
	return &CustomerRepo{DB: app.InitDb()}
}

func (repo *CustomerRepo) GetAllCustomers(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	limit := c.DefaultQuery("limit", "10")
	limitNum, _ := strconv.Atoi(limit)
	search := c.DefaultQuery("search", "")
	search = strings.ToLower(search)

	var customers []models.Customer
	query := repo.DB.Model(&models.Customer{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var count int64
	query.Count(&count)

	query = query.Offset((pageNum - 1) * limitNum).Limit(limitNum).Find(&customers)

	if query.Error != nil {
		errorMsg := query.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = map[string]interface{}{
		"customers": customers,
		"count":     count,
	}
	c.JSON(http.StatusOK, res)
}

func (repo *CustomerRepo) GetCustomerDetail(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	customerID := c.Param("id")

	var customer models.Customer
	result := repo.DB.First(&customer, customerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Customer not found"
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

	res.Data = customer
	c.JSON(http.StatusOK, res)
}

func (repo *CustomerRepo) InsertCustomer(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := ReqCustomer{}
	err := c.BindJSON(&req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorMsg := validationErrors[0].Field() + " not valid"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	req.Gender = strings.ToLower(req.Gender)
	if req.Gender != "male" && req.Gender != "female" {
		errorMsg := "Invalid gender"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	var count int64
	repo.DB.Model(&models.Customer{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		errorMsg := "Email already exists"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	customer := &models.Customer{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
	}
	result := repo.DB.Create(&customer)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = customer
	c.JSON(http.StatusCreated, res)
}

func (repo *CustomerRepo) UpdateCustomer(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := ReqUpdateCustomer{}
	err := c.BindJSON(&req)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorMsg := validationErrors[0].Field() + " not valid"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	req.Gender = strings.ToLower(req.Gender)
	if req.Gender != "male" && req.Gender != "female" {
		errorMsg := "Invalid gender"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	// Get Customer ID from URL parameter
	customerID := c.Param("id")

	// Query Customer by ID
	var customer models.Customer
	result := repo.DB.First(&customer, customerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Customer not found"
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

	// Update Customer fields
	customer.Name = req.Name
	customer.PhoneNumber = req.PhoneNumber
	customer.Gender = req.Gender

	result = repo.DB.Save(&customer)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = customer
	c.JSON(http.StatusOK, res)
}

func (repo *CustomerRepo) DeleteCustomer(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}

	customerID := c.Param("id")

	var customer models.Customer
	result := repo.DB.First(&customer, customerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			errorMsg := "Customer not found"
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

	// Delete customer's orders
	result = repo.DB.Where("customer_id = ?", customerID).Delete(&models.Order{})
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Delete customer
	result = repo.DB.Delete(&customer)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = "Customer deleted successfully"
	c.JSON(http.StatusOK, res)
}
