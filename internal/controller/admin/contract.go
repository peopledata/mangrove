package admin

import (
	"fmt"
	"mangrove/internal/controller"
	"mangrove/internal/logic"
	"mangrove/pkg/contracts"
	"strconv"
	"sync"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ContractRecordsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandDetail Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	controller.ResponseOk(c, gin.H{
		"page":     page,
		"pageSize": pageSize,
		"records":  logic.GetContractRecordsByDemandId(demandId, page, pageSize),
		"total":    logic.TotalContractRecordsByDemandId(demandId),
	})
}

// DemandContractStatusCron 需求发布后合约部署的状态更新定时器
func DemandContractStatusCron() {
	// 1. 获取所有发布中的需求（合约）
	demands := logic.GetAllPublishingDemands()
	zap.L().Debug("Demand cron job starting", zap.String("publishing", fmt.Sprintf("%d", len(demands))))

	// 2. 根据部署 tx 查询是否部署成功
	network := viper.GetString("nft.network")
	apiKey := viper.GetString("nft.infura_api_key")
	client, err := contracts.Client(network, apiKey)
	if err != nil {
		zap.L().Error("Demand contract status check cron error", zap.String("reason", "init ethclient error"), zap.Error(err))
		return
	}

	// 使用goroutine去分别执行cron worker
	var wg sync.WaitGroup
	for idx := range demands {
		wg.Add(1)
		// 执行每个需求的任务
		go logic.DemandStatusCronWorker(client, &demands[idx], &wg)
	}
	wg.Wait()
}

// DemandContractRecordsCron 需求发布后获取合约中的资产数据
func DemandContractRecordsCron() {
	// 1. 获取所有已经发布的需求（有合约地址了）
	demands := logic.GetAllPublishedDemands()
	zap.L().Debug("Demand contract records cron job starting", zap.String("published", fmt.Sprintf("%d", len(demands))))

	// 2. 根据部署 tx 查询是否部署成功
	network := viper.GetString("nft.network")
	apiKey := viper.GetString("nft.infura_api_key")
	client, err := contracts.Client(network, apiKey)
	if err != nil {
		zap.L().Error("Demand contract records cron error", zap.String("reason", "init ethclient error"), zap.Error(err))
		return
	}

	etherscanApiKey := viper.GetString("nft.etherscan_api_key")
	// 使用goroutine去分别执行cron worker
	var wg sync.WaitGroup
	for idx := range demands {
		wg.Add(1)
		// 执行每个需求的任务：分别去获取当前合约下
		go logic.DemandContractRecordsCronWorker(etherscanApiKey, client, &demands[idx], &wg)
	}
	wg.Wait()
}
