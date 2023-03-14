package logic

import (
	"fmt"
	"math/big"
	"net/url"
	"patronus/internal/dao/mysql"
	"patronus/internal/models"
	"patronus/internal/schema"
	"patronus/pkg/contracts"
	"regexp"
	"sync"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"
)

func GetContractRecordsByDemandId(demandId int64, page, pageSize int) []schema.ContractRecordResp {
	contractRecords := mysql.GetContractRecordsByDemandId(demandId, page, pageSize)
	var records []schema.ContractRecordResp
	for idx := range contractRecords {
		record := contractRecords[idx]
		// dic -> did:ppld:zQmSB8zjpigGkNNjSjLmxTgut2bQoSJ9iyW8Patv3owYnBU%40http%3A%2F%2Fppldid.peopledata.org.cn:3000
		decodedStr, _ := url.QueryUnescape(record.Did)
		r := regexp.MustCompile(`did:ppld:(.*?)@`)
		match := r.FindStringSubmatch(decodedStr)
		if len(match) > 1 {
			records = append(records, schema.ContractRecordResp{
				ID:       record.ID,
				Did:      match[1],
				SignTime: record.SignTime,
			})
		}
	}
	return records
}

func TotalContractRecordsByDemandId(demandId int64) int64 {
	return mysql.GetContractRecordsByDemandIdCount(demandId)
}

func Subscribe(client *ethclient.Client, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get a contract instance
	contractAddr := common.HexToAddress(demand.ContractAddr) // 合约地址
	instanceApi, err := contracts.NewApi(contractAddr, client)
	if err != nil {
		zap.L().Error("Demand contract events subscribe error", zap.String("reason", "get eth contract instance error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("contract_tx", demand.ContractTx), zap.Error(err))
		return
	}

	// 获取合约地址下的NFT数量
	nftCount, err := instanceApi.BalanceOf(nil, contractAddr)
	if err != nil {
		panic(err)
	}
	// 输出NFT数量
	fmt.Println("NFT数量:", nftCount)

	// Loop over all token IDs and retrieve the owner of each token
	for i := big.NewInt(0); i.Cmp(nftCount) < 0; {
		i = i.Add(i, big.NewInt(1))

		// 获取已经授权给合约地址的NFT所有者地址
		approvedAddr, err := instanceApi.GetApproved(nil, i)
		if err != nil {
			panic(err)
		}

		// 获取NFT所有者地址
		ownerAddr, err := instanceApi.OwnerOf(nil, i)
		if err != nil {
			panic(err)
		}

		// 输出NFT所有者地址和授权地址
		fmt.Println("NFT所有者地址:", ownerAddr)
		fmt.Println("授权地址:", approvedAddr)
	}

}
