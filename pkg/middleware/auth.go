package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

func AuthenticationMiddleware(jwtPublicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.ParseToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Fetch the public JWK from the Cognito endpoint
		keySet, err := fetchPublicJWTKey(context.Background(), jwtPublicKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch public JWK"})
			c.Abort()
			return
		}

		// Verify the Token
		validToken, err := verifyToken(token, keySet)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		if !validToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// Extract the subject (sub) claim from the token
		claims, ok := validToken.Claims.(jwt.MapClaims)
		if !ok || !validToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token: Subject not found"})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("sub", sub)
		c.Next()
	}
}

func verifyToken(tokenString string, keySet jwk.Set) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, errors.New("cannot look up kid header")
		}
		var publickey interface{}
		if err := keys.Raw(&publickey); err != nil {
			return nil, err
		}
		return publickey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func fetchPublicJWTKey(ctx context.Context, link string) (jwk.Set, error) {
	keySet, err := jwk.Fetch(ctx, link)
	if err != nil {
		return nil, err
	}
	return keySet, nil
}
