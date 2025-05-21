package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thanavatC/auth-service-go/controller"
)

func AuthRouter(authController *controller.AuthController, r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}
}
