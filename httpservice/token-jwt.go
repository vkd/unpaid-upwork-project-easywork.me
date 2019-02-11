package httpservice

import (
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserTokenWithClaims struct {
	Id      string
	Email   string
	Role    string
	Expires int64
}

func VerifyTokenString(tokenString string) (*UserTokenWithClaims, *TrackerError) {
	var userTokenWithClaims *UserTokenWithClaims
	pureTokenString := strings.TrimPrefix(tokenString, "Bearer ")

	var jwtParser = jwt.Parser{UseJSONNumber: true}
	token, jwtParseError := jwtParser.Parse(pureTokenString, TokenFunc)

	if validationError, ok := jwtParseError.(*jwt.ValidationError); ok {
		if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Println("That's not even a token")
			return nil, &JwtTokenParseError
		} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, &JwtTokenExpiredError
		} else {
			log.Println("Couldn't handle this token:", jwtParseError)
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if claims["jti"] == nil || claims["sub"] == nil || claims["rol"] == nil {
			return nil, &JwtTokenParseError
		}

		stClaims := jwt.StandardClaims{Subject: claims["sub"].(string), Id: claims["jti"].(string)}
		userTokenWithClaims = &UserTokenWithClaims{Email: stClaims.Subject, Expires: stClaims.ExpiresAt, Id: stClaims.Id, Role: claims["rol"].(string)}

		return userTokenWithClaims, nil

	} else {
		return nil, &JwtTokenParseError
	}
}

func TokenFunc(token *jwt.Token) (interface{}, error) {
	return []byte("karamba"), nil
}
