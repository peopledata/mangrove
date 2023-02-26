package controller

import (
	"errors"
	"patronus/internal/models"

	"github.com/gin-gonic/gin"
)

const CtxUserKey = "user"
const CtxEthKey = "ethClient"

var ErrUserNotLogin = errors.New("用户未登录")

func getCurrentUser(c *gin.Context) (*models.User, error) {
	userVal, ok := c.Get(CtxUserKey)
	if !ok {
		return nil, ErrUserNotLogin
	}
	user, ok := userVal.(*models.User)
	if !ok {
		return nil, ErrUserNotLogin
	}
	return user, nil
}
