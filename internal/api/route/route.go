package route

import (
	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/api/controller"
	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

func InitRoutes(router utils.RouterWithLogger, idpAdapter idp.CognitoAdapter) {
	userController := controller.NewUserController(idpAdapter, router.Logger)
	//
	user := router.Router.Group("/api/v2/user")
	{
		user.POST("/signup", userController.SignUp)
	}
}
