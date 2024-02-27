package route

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/api/controller"
)

func InitRoutes(router *gin.Engine, idpAdapter idp.CognitoAdapter) {
	userController := controller.NewUserController(idpAdapter)
	//
	user := router.Group("/api/v2/user")
	{
		user.POST("/signup", userController.SignUp)
	}
}
