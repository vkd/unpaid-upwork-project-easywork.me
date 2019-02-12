package httpservice

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
)

type UserTokenWithClaims struct {
	jwt.StandardClaims

	Role      models.Role `json:"rol,omitempty"`
	FirstName string      `json:"fir,omitempty"`
	LastName  string      `json:"las,omitempty"`
}

type claimCreator struct {
	secret []byte
}

type ClaimCreator interface {
	CreateClaim(*models.User) (string, error)
}

func (c *claimCreator) CreateClaim(u *models.User) (string, error) {
	var claim UserTokenWithClaims
	claim.ExpiresAt = time.Now().AddDate(0, 0, 7).Unix()
	claim.Id = string(u.ID)
	claim.Subject = u.Email
	claim.Role = u.Role
	claim.FirstName = u.FirstName
	claim.LastName = u.LastName

	str, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(c.secret)
	return str, err
}

func VerifyTokenString(tokenString string, secret []byte) (*models.User, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &UserTokenWithClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, &models.JwtTokenParseError
			} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, &models.JwtTokenExpiredError
			}
		}
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
	user.Role = models.Role(c.Role)

	return &user, nil
}

func TokenFunc(token *jwt.Token) (interface{}, error) {
	return []byte("karamba"), nil
}
