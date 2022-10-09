package main

import (
	"go_crowdfund/auth"
	"go_crowdfund/campaign"
	"go_crowdfund/handler"
	"go_crowdfund/helper"
	"go_crowdfund/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandle := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.POST("/campaign", authMiddleware(authService, userService), campaignHandle.CreateCampaign)
	api.PUT("/campaign/:id", authMiddleware(authService, userService), campaignHandle.UpdateCampaign)
	api.POST("/campaign-image", authMiddleware(authService, userService), campaignHandle.UploadImage)

	api.GET("/campaigns", campaignHandle.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandle.GetCampaign)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(http.StatusUnauthorized, "Unauthorized", "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		validateToken, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, "Unauthorized", "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := validateToken.Claims.(jwt.MapClaims)

		if !ok || !validateToken.Valid {
			response := helper.APIResponse(http.StatusUnauthorized, "Unauthorized", "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, "Unauthorized", "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
