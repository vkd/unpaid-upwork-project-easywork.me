package httpservice

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
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

func ObjectIDParam(c *gin.Context, name string) (bson.ObjectId, error) {
	s := c.Param(name)

	d, err := hex.DecodeString(s)
	if err != nil || len(d) != 12 {
		return "", fmt.Errorf("invalid input to ObjectIdHex: %q", s)
	}
	return bson.ObjectId(d), nil
}
