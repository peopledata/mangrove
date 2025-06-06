package admin

import (
	"errors"
	"mangrove/internal/controller"
	"mangrove/internal/dao/mysql"
	"mangrove/internal/logic"
	"mangrove/internal/schema"
	"mangrove/pkg/converter"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func DemandListHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	q := c.DefaultQuery("q", "")

	controller.ResponseOk(c, gin.H{
		"page":     page,
		"pageSize": pageSize,
		"demands":  logic.ListPagerDemands(q, page, pageSize),
		"total":    logic.TotalDemands(q),
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

	// 设置成至少是当前时间一小时后
	now := time.Now().Add(time.Hour)
	if dcr.ValidAt.Unix() < now.Unix() {
		zap.L().Error("DemandCreate Handler valid at lt now", zap.String("name", dcr.Name))
		controller.ResponseErr(c, controller.CodeDemandInvalidAtErr)
		return
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		zap.L().Error("DemandCreate Handler get current user error", zap.String("name", dcr.Name), zap.Error(err))
		controller.ResponseErr(c, controller.CodeNeedAuth)
		return
	}
	//	2. 业务处理
	if err := logic.CreateDemand(&dcr, user); err != nil {
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

	// 2. 参数校验
	var dur schema.DemandUpdateReq
	dur.DemandId = demandIdStr
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

	// 设置成至少是当前时间一小时后
	now := time.Now().Add(time.Hour)
	if dur.ValidAt.Unix() < now.Unix() {
		zap.L().Error("DemandUpdate Handler valid at lt now", zap.String("name", dur.Name))
		controller.ResponseErr(c, controller.CodeDemandInvalidAtErr)
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
	marketPlaceHost := viper.GetString("marketplace.host")
	marketPlaceApiKey := viper.GetString("marketplace.api_key")
	if err := logic.PublishDemand(demandId, marketPlaceHost, marketPlaceApiKey, ethClient); err != nil {
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
