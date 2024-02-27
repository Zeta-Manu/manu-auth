package idp

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

func handleCognitoError(err error) error {
	var invalidPasswordErr *types.InvalidPasswordException
	var invalidParameterErr *types.InvalidParameterException
	var usernameExistsErr *types.UsernameExistsException
	var notAuthorizedErr *types.NotAuthorizedException
	var userNotFoundErr *types.UserNotFoundException
	var userNotConfirmErr *types.UserNotConfirmedException
	var aliasExistErr *types.AliasExistsException

	switch {
	case errors.As(err, &invalidPasswordErr):
		return &utils.CustomError{
			Message: "Invalid password",
			Status:  http.StatusBadRequest,
		}
	case errors.As(err, &invalidParameterErr):
		return &utils.CustomError{
			Message: "Invalid parameter",
			Status:  http.StatusBadRequest,
		}
	case errors.As(err, &usernameExistsErr):
		return &utils.CustomError{
			Message: "Username already exists",
			Status:  http.StatusConflict,
		}
	case errors.As(err, &notAuthorizedErr):
		return &utils.CustomError{
			Message: "Not Authorized",
			Status:  http.StatusUnauthorized,
		}
	case errors.As(err, &userNotFoundErr):
		return &utils.CustomError{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
	case errors.As(err, &userNotConfirmErr):
		return &utils.CustomError{
			Message: "User not confirm",
			Status:  http.StatusForbidden,
		}
	case errors.As(err, &aliasExistErr):
		return &utils.CustomError{
			Message: "Alias exists",
			Status:  http.StatusConflict,
		}
	default:
		return &utils.CustomError{
			Message: "Internal error",
			Status:  http.StatusInternalServerError,
		}
	}
}
