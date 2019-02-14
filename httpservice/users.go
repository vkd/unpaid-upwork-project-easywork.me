package httpservice

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func usersGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		us, err := db.UsersGetPublic(c)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, us)
	}
}

func userGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		u, err := db.UserGet(c, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, u)
	}
}

func userCreateHandler(db *storage.Storage, cc ClaimCreator) gin.HandlerFunc {
	type UserData struct {
		ID        models.UserID `json:"id" db:"id"`
		Email     string        `json:"email" db:"email"`
		Password  string        `json:"password" db:"password"`
		UserType  models.Role   `json:"user_type" db:"user_type"`
		FirstName string        `json:"first_name" db:"first_name"`
		LastName  string        `json:"last_name" db:"last_name"`
	}
	return func(c *gin.Context) {
		var u UserData
		err := c.ShouldBindJSON(&u)
		if err != nil {
			apiError(c, http.StatusBadRequest, &models.JsonDecodeError)
			return
		}

		if u.Email == "" || u.Password == "" {
			apiError(c, http.StatusBadRequest, &models.UserEmailOrPasswordEmpty)
			return
		}

		if !IsUsernameValid(u.ID) {
			apiError(c, http.StatusBadRequest, &models.WrongUsername)
			return
		}

		if u.FirstName == "" {
			apiError(c, http.StatusBadRequest, &models.EmptyFirstName)
			return
		}
		if u.LastName == "" {
			apiError(c, http.StatusBadRequest, &models.EmptyLastName)
			return
		}

		_, err = db.UserGetByEmail(c, u.Email)
		if err == nil {
			apiError(c, http.StatusUnprocessableEntity, &models.UserEmailExists)
			return
		}
		if !storage.IsNotFound(err) {
			apiError(c, http.StatusInternalServerError, errors.Wrapf(err, "error on find user by email (email: %v)", u.Email))
			return
		}

		user := models.UserPassword{
			User: models.User{
				ID:        u.ID,
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Role:      u.UserType,
			},
			Password: u.Password,
		}
		out, err := db.UserCreate(c, &user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		token, err := cc.CreateClaim(out)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, struct {
			models.User
			Token string `json:"token"`
		}{
			User:  *out,
			Token: token,
		})
	}
}

func userDeleteHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		err := db.UserDelete(c, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, UserDeleted)
	}
}

func IsUsernameValid(username models.UserID) bool {
	r, _ := regexp.Compile(`^[a-z0-9]+$`)
	return r.MatchString(string(username))
}
