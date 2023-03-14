package mysql

import "errors"

var (
	ErrUserExist       = errors.New("用户已存在")
	ErrUserNotExist    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")

	ErrDemandNotExist         = errors.New("需求不存在")
	ErrDemandStatusNotInit    = errors.New("当前需求不是草稿状态，不能进行发布")
	ErrContractRecordNotExist = errors.New("签约记录不存在")
)
