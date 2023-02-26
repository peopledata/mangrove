package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	TokenExpireDuration = time.Hour * 2
)

var SignSecret = []byte("PeopleData")

var keyFunc = func(token *jwt.Token) (interface{}, error) {
	return SignSecret, nil
}

type MyClaims struct {
	ID       uint   `json:"id"`
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenToken(id uint, userId int64, username string) (aToken, rToken string, err error) {
	// 创建一个自己的声明
	c := MyClaims{
		ID:       id,
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "PeopleData",                               // 签发人
		},
	}
	//	使用指定的签名方法创建签名对象
	//	使用指定的secret签名并获得完整的编码后的token字符串
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(SignSecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
		Issuer:    "PeopleData",
	}).SignedString(SignSecret)

	return
}

func ParseToken(tokenStr string) (*MyClaims, error) {
	// 解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenStr, mc, keyFunc)
	if err != nil {
		return nil, err
	}
	// 校验 token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("token无效")
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	//	refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧的access token中解析出claims数据
	var mc = new(MyClaims)
	_, err = jwt.ParseWithClaims(aToken, mc, keyFunc)
	v, _ := err.(*jwt.ValidationError)
	//	当access token是过期错误并 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(mc.ID, mc.UserId, mc.Username)
	}
	return
}
