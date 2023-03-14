package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"patronus/internal/controller"
	"patronus/internal/logic"
	"patronus/internal/models"
	"patronus/pkg/contracts"
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
	alchemyApiKey := viper.GetString("nft.alchemy_api_key")
	// http://localhost:8545
	client, err := contracts.Client(alchemyApiKey)
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

// DemandContractSubscribeWorker 监听需求合约事件
// TODO：订阅方式有问题
func DemandContractSubscribeWorker(etherscanApiKey string) {
	// 1. 获取所有已经发布的需求（有合约地址了）
	demands := logic.GetAllPublishedDemands()
	zap.L().Debug("Demand contract subscribe worker starting", zap.String("published", fmt.Sprintf("%d", len(demands))))

	// 2. 根据部署 tx 查询是否部署成功
	//alchemyApiKey := viper.GetString("nft.alchemy_api_key")
	// http://localhost:8545
	//client, err := contracts.Client(alchemyApiKey)
	//if err != nil {
	//	zap.L().Error("Demand contract records cron error", zap.String("reason", "init ethclient error"), zap.Error(err))
	//	return
	//}

	// 使用goroutine去分别执行cron worker
	var wg sync.WaitGroup
	for idx := range demands {
		wg.Add(1)
		// 执行每个需求的任务：分别去获取当前合约下
		//go logic.Subscribe(client, &demands[idx], &wg)
		//go logic.SubscribeContract(client, &demands[idx], &wg)
		go ContractTransferEvent(etherscanApiKey, &demands[idx], &wg)
	}
	wg.Wait()
}

func ContractTransferEvent(etherscanApiKey string, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	// 查询该NFT的交易历史记录
	url := fmt.Sprintf("https://api-goerli.etherscan.io/api?module=account&action=tokennfttx&contractaddress=%s&tokenid=%d&page=1&offset=100&sort=desc&apikey=%s",
		demand.ContractAddr, 1, etherscanApiKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	if result["status"].(string) == "1" {
		transactions := result["result"].([]interface{})
		fmt.Println("查询到的NFT交易历史记录：")
		for _, tx := range transactions {
			txHash := tx.(map[string]interface{})["hash"].(string)
			blockNumber := tx.(map[string]interface{})["blockNumber"].(string)
			timeStamp := tx.(map[string]interface{})["timeStamp"].(string)
			from := tx.(map[string]interface{})["from"].(string)
			to := tx.(map[string]interface{})["to"].(string)
			tokenID := tx.(map[string]interface{})["tokenID"].(string)

			fmt.Printf("TokenID: %s\n交易哈希：%s\n区块高度：%s\n交易时间戳：%s\n发送地址：%s\n接收地址：%s\n\n",
				tokenID, txHash, blockNumber, timeStamp, from, to)

			// 如果是发给合约地址的，则表示授权给合约地址的
			if from != "0x0000000000000000000000000000000000000000" && to == demand.ContractAddr {

			}

		}
	} else {
		fmt.Println("查询失败：", result["message"].(string))
	}

	// 查询包含该NFT的所有交易
	//url = fmt.Sprintf("https://api-goerli.etherscan.io/api?module=account&action=txlist&address=%s&page=1&offset=100&sort=desc&apikey=%s",
	//	demand.ContractAddr, etherscanApiKey)
	//fmt.Println(url)
	//resp, err = http.Get(url)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//err = json.NewDecoder(resp.Body).Decode(&result)
	//if err != nil {
	//	panic(err)
	//}
	//
	//if result["status"].(string) == "1" {
	//	transactions := result["result"].([]interface{})
	//	fmt.Println("查询到的包含该NFT的交易：")
	//	for _, tx := range transactions {
	//		txHash := tx.(map[string]interface{})["hash"].(string)
	//		blockNumber := tx.(map[string]interface{})["blockNumber"].(string)
	//		timeStamp := tx.(map[string]interface{})["timeStamp"].(string)
	//		value := tx.(map[string]interface{})["value"].(string)
	//		input := tx.(map[string]interface{})["input"].(string)
	//
	//		// 判断该交易是否包含该NFT
	//		if strings.Contains(input, "000000000000000000000000"+"1") {
	//			from := tx.(map[string]interface{})["from"].(string)
	//			to := tx.(map[string]interface{})["to"].(string)
	//
	//			fmt.Printf("交易哈希：%s\n区块高度：%s\n交易时间戳：%s\n发送地址：%s\n接收地址：%s\n交易金额：%s\n\n",
	//				txHash, blockNumber, timeStamp, from, to, value)
	//		}
	//	}
	//} else {
	//	fmt.Println("查询失败：", result["message"].(string))
	//}
}
