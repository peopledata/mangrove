package admin

import (
	"errors"
	"patronus/internal/controller"
	"patronus/internal/dao/mysql"
	"patronus/internal/logic"
	"patronus/internal/schema"
	"patronus/pkg/converter"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func DemandListHandler(c *gin.Context) {
	// todo：分页
	controller.ResponseOk(c, gin.H{
		"demands": logic.ListDemands(),
		"total":   logic.TotalDemands(),
	})

}

// DemandDetailHandler 获取详细信息，详情页面展示
func DemandDetailHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandDetail Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	data, err := logic.GetDemandDetail(demandId)
	if err != nil {
		zap.L().Error("DemandDetail Handler get detail error", zap.String("id", demandIdStr), zap.Error(err))
		if errors.Is(err, mysql.ErrDemandNotExist) {
			controller.ResponseErr(c, controller.CodeDemandNotExist)
		} else {
			controller.ResponseErr(c, controller.CodeDemandCreateErr)
		}
		return
	}

	// 3. 返回Response
	controller.ResponseOk(c, data)
}

// DemandInfoHandler 获取基本信息，主要用于更新需求的数据回填
func DemandInfoHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandInfo Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	data, err := logic.GetDemandInfo(demandId)
	if err != nil {
		zap.L().Error("DemandInfo Handler get detail error", zap.String("id", demandIdStr), zap.Error(err))
		if errors.Is(err, mysql.ErrDemandNotExist) {
			controller.ResponseErr(c, controller.CodeDemandNotExist)
		} else {
			controller.ResponseErr(c, controller.CodeDemandCreateErr)
		}
		return
	}

	// 3. 返回Response
	controller.ResponseOk(c, data)
}

func DemandCreateHandler(c *gin.Context) {
	//	1.参数校验
	var dcr schema.DemandCreateReq
	if err := c.ShouldBindJSON(&dcr); err != nil {
		zap.L().Error("DemandCreate Handler with invalid param", zap.Error(err))
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
	if err := logic.CreateDemand(&dcr); err != nil {
		zap.L().Error("DemandCreate Handler logic handle error", zap.String("name", dcr.Name), zap.Error(err))
		controller.ResponseErr(c, controller.CodeDemandCreateErr)
		return
	}

	//	3.返回响应
	controller.ResponseOk(c, nil)
}

func DemandUpdateHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandUpdate Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 参数校验
	var dur schema.DemandUpdateReq
	dur.DemandId = demandId
	if err := c.ShouldBindJSON(&dur); err != nil {
		zap.L().Error("DemandUpdate Handler with invalid param", zap.Error(err))
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

	//	3. 业务处理
	if err := logic.UpdateDemand(&dur); err != nil {
		zap.L().Error("DemandUpdate Handler logic handle error", zap.String("demand_id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeDemandUpdateErr)
		return
	}

	//	4.返回响应
	controller.ResponseOk(c, nil)

}

func DemandPublishHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandPublish Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 获取 etc client
	val, ok := c.Get(controller.CtxEthKey)
	if !ok {
		controller.ResponseErr(c, controller.CodeUnknown)
		return
	}
	ethClient := val.(*ethclient.Client)

	// 3. 执行发布逻辑
	if err := logic.PublishDemand(demandId, ethClient); err != nil {
		zap.L().Error("DemandPublish Handler logic handle error", zap.String("demand_id", demandIdStr), zap.Error(err))
		if errors.Is(err, mysql.ErrDemandStatusNotInit) {
			controller.ResponseErr(c, controller.CodeDemandStatusNotInit)
		} else {
			controller.ResponseErr(c, controller.CodeDemandPublishErr)
		}
		return
	}

	controller.ResponseOk(c, nil)
}

func DemandDeleteHandler(c *gin.Context) {
	// 1. 获取demand id
	//demandIdStr := c.Param("id")
	//demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	//if err != nil {
	//	zap.L().Error("DemandPublish Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
	//	ResponseErr(c, CodeInvalidParam)
	//	return
	//}
}
