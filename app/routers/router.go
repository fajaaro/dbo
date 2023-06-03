package routers

import (
	"github.com/fajaaro/dbo/app/controllers"
	"github.com/fajaaro/dbo/app/middlewares"
	"github.com/gin-gonic/gin"
)

type API struct {
	AuthRepo     controllers.AuthRepo
	OrderRepo    controllers.OrderRepo
	CustomerRepo controllers.CustomerRepo
}

func SetupRouter(AuthRepo controllers.AuthRepo, OrderRepo controllers.OrderRepo, CustomerRepo controllers.CustomerRepo) *gin.Engine {
	r := gin.New()
	api := API{
		AuthRepo,
		OrderRepo,
		CustomerRepo,
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(middleware.CORSMiddleware(), middleware.ClientInfo())

	authRoutes := r.Group("")
	authRoutes.POST("/api/auth/register", api.AuthRepo.Register)
	authRoutes.POST("/api/auth/login", api.AuthRepo.Login)
	authRoutes.POST("/api/auth/refresh-token", api.AuthRepo.RefreshToken)
	authRoutes.POST("/api/auth/match-token", api.AuthRepo.MatchToken)

	orderRoutes := r.Group("")
	orderRoutes.Use(middlewares.JWT())
	orderRoutes.GET("/api/orders", api.OrderRepo.GetAllOrders)
	orderRoutes.GET("/api/orders/:id", api.OrderRepo.GetOrderDetail)
	orderRoutes.POST("/api/orders", api.OrderRepo.InsertOrder)
	orderRoutes.PUT("/api/orders/:id", api.OrderRepo.UpdateOrder)
	orderRoutes.DELETE("/api/orders/:id", api.OrderRepo.DeleteOrder)

	customerRoutes := r.Group("")
	customerRoutes.Use(middlewares.JWT())
	customerRoutes.GET("/api/customers", api.CustomerRepo.GetAllCustomers)
	customerRoutes.GET("/api/customers/:id", api.CustomerRepo.GetCustomerDetail)
	customerRoutes.POST("/api/customers", api.CustomerRepo.InsertCustomer)
	customerRoutes.PUT("/api/customers/:id", api.CustomerRepo.UpdateCustomer)
	customerRoutes.DELETE("/api/customers/:id", api.CustomerRepo.DeleteCustomer)

	return r
}
