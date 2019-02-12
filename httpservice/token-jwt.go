package httpservice

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
)

type UserTokenWithClaims struct {
	jwt.StandardClaims

	Rol models.UserRole `json:"rol,omitempty"`
}

// func CreateTokenString() (string, error) {
// 	jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
// }

func VerifyTokenString(tokenString string, secret []byte) (*models.User, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &UserTokenWithClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if validationError, ok := err.(*jwt.ValidationError); ok {
		if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, &models.JwtTokenParseError
		} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, &models.JwtTokenExpiredError
		}
	}
	if err != nil {
		return nil, err
	}

	var c *UserTokenWithClaims
	var ok bool
	if c, ok = token.Claims.(*UserTokenWithClaims); !ok {
		return nil, errors.Errorf("wrong token claim type")
	}
	if c == nil {
		return nil, errors.Errorf("wrong user claim (is nil)")
	}

	var user models.User
	user.ID = models.UserID(c.StandardClaims.Id)
	user.Email = c.StandardClaims.Subject
	user.Role = models.UserRole(c.Rol)

	return &user, nil
}

func TokenFunc(token *jwt.Token) (interface{}, error) {
	return []byte("karamba"), nil
}
