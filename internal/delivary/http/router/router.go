package router

import (
	"net/http"

	"github.com/alirezaghasemi/user-manager/internal/delivary/http/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler handler.UserHandler) *gin.Engine {
	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})

	baseRouter := router.Group("/api/v1")
	userRouter := baseRouter.Group("/user")

	// Create User
	userRouter.POST("", userHandler.Create)

	// Get User
	userRouter.GET("/:id", userHandler.FindByID)

	// Update User
	userRouter.PATCH("/:id", userHandler.Update)

	// Get All Users
	userRouter.GET("", userHandler.FindAll)

	// Delete User
	userRouter.DELETE("/:id", userHandler.Delete)
	return router
}
