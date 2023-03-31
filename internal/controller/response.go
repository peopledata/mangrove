package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeSignUpErr
	CodeUserNotExist
	CodeInvalidPassword
	CodeLoginErr

	CodeNeedAuth
	CodeInvalidToken
	CodeInvalidEthClient

	CodeDemandCreateErr
	CodeDemandInvalidAtErr
	CodeDemandDetailErr
	CodeDemandUpdateErr
	CodeDemandNotExist
	CodeDemandStatusNotInit
	CodeDemandPublishErr

	CodeTaskCreateErr
	CodeTaskListErr

	CodeAlgoRecordsErr

	CodeServerBusy
	CodeUnknown
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeSignUpErr:       "用户创建失败",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeLoginErr:        "登录失败",

	CodeNeedAuth:     "需要登录认证",
	CodeInvalidToken: "无效的认证Token",

	CodeInvalidEthClient: "无效的Eth客户端",

	CodeDemandCreateErr:     "新增需求失败",
	CodeDemandInvalidAtErr:  "有效期错误",
	CodeDemandDetailErr:     "获取需求失败",
	CodeDemandUpdateErr:     "更新需求失败",
	CodeDemandNotExist:      "需求不存在",
	CodeDemandStatusNotInit: "当前需求不是草稿状态，不能进行发布",
	CodeDemandPublishErr:    "需求发布失败",

	CodeTaskCreateErr: "任务创建失败",
	CodeTaskListErr:   "任务列表获取失败",

	CodeAlgoRecordsErr: "获取算法记录失败",

	CodeServerBusy: "服务繁忙",
	CodeUnknown:    "未知错误",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeUnknown]
	}
	return msg
}

/*
{
	"code": 10001, // 程序中的错误码
	"msg": "提示信息",
    "data": {}  // 存放数据
}
*/

type ResponseData struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseErr(c *gin.Context, code ResCode) {
	rd := ResponseData{
		Code: int(code),
		Msg:  code.Msg(),
	}
	c.JSON(http.StatusOK, &rd)
}

func ResponseErrWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := ResponseData{
		Code: int(code),
		Msg:  msg,
	}
	c.JSON(http.StatusOK, &rd)
}

func ResponseOk(c *gin.Context, data interface{}) {
	rd := ResponseData{
		Code: int(CodeSuccess),
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, &rd)
}

func ResponseOkWithMsg(c *gin.Context, data interface{}, msg interface{}) {
	rd := ResponseData{
		Code: int(CodeSuccess),
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, &rd)
}
