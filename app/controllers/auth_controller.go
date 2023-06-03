package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/fajaaro/dbo/app"
	"github.com/fajaaro/dbo/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

type ReqAuth struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func AuthController() *AuthRepo {
	return &AuthRepo{DB: app.InitDb()}
}

func ValidateAccessToken(accessToken string, db *gorm.DB) (*models.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return nil, errors.New("Invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid access token")
	}

	userEmail, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("Invalid user email in token claims")
	}

	var user models.User
	result := db.Where("email = ?", userEmail).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))

func createToken(tokenType string, exp time.Time, email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_type": tokenType,
		"exp":        exp.Unix(),
		"sub":        email,
	})

	tokenStr, _ := token.SignedString(SECRET_KEY)

	return tokenStr
}

func (repo *AuthRepo) Register(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := ReqAuth{}
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

	var count int64
	repo.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		errorMsg := "Email already exists"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		res.Success = false
		errorMsg := "Failed to hash password"
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		c.Abort()
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	result := repo.DB.Create(&user)
	if result.Error != nil {
		errorMsg := result.Error.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusInternalServerError, res)
		c.Abort()
		return
	}

	res.Data = map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}
	c.JSON(http.StatusCreated, res)
}

func (repo *AuthRepo) Login(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := ReqAuth{}
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

	var user models.User
	result := repo.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		errorMsg := "Invalid credentials"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		errorMsg := "Invalid credentials"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	accessTokenExp := time.Now().Add(15 * time.Minute)
	accessToken := createToken("access_token", accessTokenExp, user.Email)

	refreshTokenExp := time.Now().Add(30 * 24 * time.Hour)
	refreshToken := createToken("refresh_token", refreshTokenExp, user.Email)

	data := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	res.Data = data
	c.JSON(http.StatusOK, res)
}

func (repo *AuthRepo) RefreshToken(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := map[string]string{}
	err := c.BindJSON(&req)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	refreshToken := req["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		errorMsg := "Invalid refresh token"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errorMsg := "Expired refresh token"
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	userEmail := claims["sub"].(string)

	accessTokenExp := time.Now().Add(15 * time.Minute)
	accessToken := createToken("access_token", accessTokenExp, userEmail)

	data := gin.H{
		"access_token": accessToken,
	}
	res.Data = data
	c.JSON(http.StatusOK, res)
}

func (repo *AuthRepo) MatchToken(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	res := models.JsonResponse{Success: true}
	req := map[string]string{}
	err := c.BindJSON(&req)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	accessToken := req["access_token"]

	user, err := ValidateAccessToken(accessToken, repo.DB)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg

		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	res.Data = map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}
	c.JSON(http.StatusOK, res)
}
