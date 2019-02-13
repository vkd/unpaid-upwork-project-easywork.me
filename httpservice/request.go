package httpservice

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/pkg/errors"
)

// IntParam - get param from context and convert
func IntParam(c *gin.Context, name string) (int, error) {
	p := c.Param(name)
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, errors.Wrapf(err, "param %q is not int (param: %q)", name, p)
	}
	return i, nil
}

func ObjectIDParam(c *gin.Context, name string) (primitive.ObjectID, error) {
	s := c.Param(name)
	return primitive.ObjectIDFromHex(s)
}
