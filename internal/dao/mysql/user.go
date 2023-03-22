package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"mangrove/internal/models"

	"gorm.io/gorm"
)

const secret = "PeopleData"

func CheckUserExist(username string) error {
	var count int64
	if err := db.Model(&models.User{}).Where("username=?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExist
	}
	return nil
}

func InsertUser(user *models.User) error {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	return db.Create(user).Error
}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(user *models.User) error {
	originPassword := user.Password
	err := db.Where("username = ?", user.Username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return ErrUserNotExist
	} else if err != nil {
		return err
	} else if encryptPassword(originPassword) != user.Password {
		return ErrInvalidPassword
	}
	return nil
}
