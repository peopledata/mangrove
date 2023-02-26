package admin

import (
	"errors"
	"fmt"
	"patronus/internal/controller"
	"patronus/internal/dao/mysql"
	"patronus/internal/logic"
	"patronus/internal/schema"
	"patronus/pkg/contracts"
	"patronus/pkg/converter"
	"strconv"
	"sync"

	"github.com/spf13/viper"

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

// DemandContractStatusCron 需求发布后合约部署的状态更新定时器
func DemandContractStatusCron() {
	// 1. 获取所有发布中的需求（合约）
	demands := logic.GetAllPublishingDemands()
	zap.L().Debug("Demand cron job starting", zap.String("publishing", fmt.Sprintf("%d", len(demands))))

	// 2. 根据部署 tx 查询是否部署成功
	alchemyApiKey := viper.GetString("nft.alchemy_api_key")
	// http://localhost:8545
	client, err := contracts.Client(alchemyApiKey)
	if err != nil {
		zap.L().Error("Demand contract status check cron error", zap.String("reason", "init ethclient error"), zap.Error(err))
		return
	}

	// 使用goroutine去分别执行cron worker
	var wg sync.WaitGroup
	for idx := range demands {
		wg.Add(1)
		// 执行每个需求的任务
		go logic.DemandCronWorker(client, &demands[idx], &wg)
	}
	wg.Wait()
}
