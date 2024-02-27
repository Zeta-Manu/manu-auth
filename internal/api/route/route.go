package route

import (
	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/api/controller"
	"github.com/Zeta-Manu/manu-auth/pkg/middleware"
	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

func InitRoutes(router utils.RouterWithLogger, idpAdapter idp.CognitoAdapter, jwtPublicKey string) {
	userController := controller.NewUserController(idpAdapter, router.Logger)
	//
	user := router.Router.Group("/api/v2/user")
	{
		user.POST("/signup", userController.SignUp)
		user.POST("/confirm", userController.ConfirmSignUp)
		user.POST("/resend-confirm", userController.ResendConfirmationCode)
		user.POST("/login", userController.LogIn)
		user.POST("/forgot-password", userController.ForgotPassword)
		user.POST("/confirm-forgot", userController.ConfirmForgotPassword)
		// route with middleware
		user.POST("/password", middleware.AuthenticationMiddleware(jwtPublicKey), userController.ChangePassword)
	}
}
