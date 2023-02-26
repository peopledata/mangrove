package admin

import (
	"errors"
	"patronus/internal/controller"
	"patronus/internal/dao/mysql"
	"patronus/internal/logic"
	"patronus/internal/models"
	"patronus/internal/schema"
	"patronus/pkg/converter"
	"patronus/pkg/jwt"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//	1.参数校验
	var sp schema.SignUpReq
	if err := c.ShouldBindJSON(&sp); err != nil {
		zap.L().Error("SignUp Handler with invalid param", zap.Error(err))
		// 判断err是不是ValidationErrors错误，如果是则翻译下错误为中文
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 不是ValidationErrors类型错误直接返回
			controller.ResponseErr(c, controller.CodeInvalidParam)
			return
		}
		controller.ResponseErrWithMsg(c, controller.CodeInvalidParam, converter.RemoveTopStruct(errs.Translate(converter.Trans))) // 翻译错误信息
		return
	}

	//	2. 业务处理
	if err := logic.SignUp(&sp); err != nil {
		zap.L().Error("SignUp Handler logic handle error", zap.String("username", sp.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrUserExist) {
			controller.ResponseErr(c, controller.CodeUserExist)
			return
		}
		controller.ResponseErr(c, controller.CodeSignUpErr)
		return
	}

	//	3.返回响应
	controller.ResponseOk(c, nil)
}

func LoginHandler(c *gin.Context) {
	//	1.获取请求参数及参数校验
	var lp schema.LoginReq
	if err := c.ShouldBindJSON(&lp); err != nil {
		zap.L().Error("Login Handler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			controller.ResponseErr(c, controller.CodeInvalidParam)
			return
		}
		controller.ResponseErrWithMsg(c, controller.CodeInvalidParam, converter.RemoveTopStruct(errs.Translate(converter.Trans))) // 翻译错误信息
		return
	}

	//	2.业务逻辑处理
	accessToken, refreshToken, err := logic.Login(&lp)
	if err != nil {
		zap.L().Error("Login Handler logic handle error", zap.String("username", lp.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			controller.ResponseErr(c, controller.CodeUserNotExist)
		} else if errors.Is(err, mysql.ErrInvalidPassword) {
			controller.ResponseErr(c, controller.CodeInvalidPassword)
		} else {
			controller.ResponseErr(c, controller.CodeLoginErr)
		}
		return
	}

	//	3.返回响应
	controller.ResponseOk(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	//	1.获取请求参数及参数校验
	var rtp schema.RefreshTokenReq
	if err := c.ShouldBindJSON(&rtp); err != nil {
		zap.L().Error("RefreshToken Handler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			controller.ResponseErr(c, controller.CodeInvalidParam)
			return
		}
		controller.ResponseErrWithMsg(c, controller.CodeInvalidParam, converter.RemoveTopStruct(errs.Translate(converter.Trans))) // 翻译错误信息
		return
	}

	// 2.刷新JWT Token
	newAccessToken, _, err := jwt.RefreshToken(rtp.AccessToken, rtp.RefreshToken)
	if err != nil {
		controller.ResponseErr(c, controller.CodeInvalidToken)
		return
	}

	//	3.返回响应
	controller.ResponseOk(c, gin.H{
		"access_token": newAccessToken,
	})

}

func GetUserHandler(c *gin.Context) {
	val, ok := c.Get(controller.CtxUserKey)
	if !ok {
		controller.ResponseErr(c, controller.CodeUnknown)
		return
	}
	user := val.(*models.User)
	controller.ResponseOk(c, &schema.UserInfoResp{
		Id:       user.UserId,
		Username: user.Username,
		Permissions: &schema.UserRole{
			Role: "admin",
		},
		Avatar: "https://bxdc-static.oss-cn-beijing.aliyuncs.com/images/1675777586153.jpg",
	})
}
