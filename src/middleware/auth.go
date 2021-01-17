package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"limecord-core/utils"
	"strings"
)

var (
	TOKEN_SECRET = []byte("epic_nexure")
)

type AuthClaims struct {
	UserID string
	IsBot bool
}

func NewAuthClaims(claims jwt.MapClaims) (AuthClaims, error) {
	var (
		UserID string
		IsBot bool
	)

	fmt.Printf("%v\n", claims)

	return AuthClaims {
		UserID,
		IsBot,
	}, nil
}

func parseJwt(tokenStr string) (*AuthClaims, error) {
	var (
		err error = nil
		token *jwt.Token
	)

	if token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return TOKEN_SECRET, nil
	}); err != nil || token == nil {
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to parse/verify token")
	}


	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("failed to verify token")
	} else if authClaims, err := NewAuthClaims(claims); err != nil {
		return nil, err
	} else {
		return &authClaims, nil
	}
}

func AuthMiddleware(ctx *gin.Context) {
	isBot := false
	authHeader := ctx.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Bot ") {
		isBot = true
		authHeader = authHeader[4:]
	}

	if claims, err := parseJwt(authHeader); err != nil || isBot != claims.IsBot {
		ctx.AbortWithStatusJSON(403, utils.GetAPIError(utils.API_UNAUTHORIZED, utils.API_UNAUTHORIZED_MESSAGE))
	} else {
		ctx.Set("auth", claims)
	}

}
