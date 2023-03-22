package logic

import (
	"mangrove/internal/dao/mysql"
	"mangrove/internal/models"
	"mangrove/internal/schema"
	"mangrove/pkg/jwt"
	"mangrove/pkg/snowflake"
)

func SignUp(sp *schema.SignUpReq) error {
	// 1. 判断用户是否存在
	if err := mysql.CheckUserExist(sp.Username); err != nil {
		return err
	}

	// 2. 生成UID
	userId := snowflake.GenID()

	//	3. 构造一个User结构体
	user := models.User{
		UserId:   userId,
		Username: sp.Username,
		Password: sp.Password,
	}
	return mysql.InsertUser(&user)
}

func Login(lp *schema.LoginReq) (string, string, error) {
	user := models.User{
		Username: lp.Username,
		Password: lp.Password,
	}
	if err := mysql.Login(&user); err != nil {
		return "", "", err
	}
	//	生成JWT Token
	return jwt.GenToken(user.ID, user.UserId, user.Username)
}
