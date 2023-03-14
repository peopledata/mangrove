package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"patronus/internal/dao/mysql"
	"patronus/internal/models"
	"patronus/internal/schema"
	"patronus/pkg/contracts"
	"patronus/pkg/ipfs"
	"patronus/pkg/snowflake"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/interface-go-ipfs-core/path"

	ipfsapi "github.com/ipfs/go-ipfs-api"
	//"github.com/ipfs/go-path"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"go.uber.org/zap"
)

func CreateDemand(dcr *schema.DemandCreateReq) error {
	// 2. 生成UID
	demandId := snowflake.GenID()

	//	3. 构造一个User结构体
	demand := models.Demand{
		DemandId:  demandId,
		Name:      dcr.Name,
		Brief:     dcr.Brief,
		ValidAt:   dcr.ValidAt,
		Status:    models.DemandStatusInit,
		Category:  dcr.Category,
		Content:   dcr.Content,
		NeedUsers: dcr.NeedUsers,
		UseTimes:  dcr.UseTimes,
		Purpose:   dcr.Purpose,
		Algorithm: dcr.Algorithm,
		Agreement: dcr.Agreement,
	}
	return mysql.InsertDemand(&demand)
}

func ListDemands() []schema.DemandListResp {
	demands := mysql.GetAllDemands()
	var demandList []schema.DemandListResp
	for idx := range demands {
		item := demands[idx]
		demandList = append(demandList, schema.DemandListResp{
			DemandId:       item.DemandId,
			Name:           item.Name,
			ValidAt:        item.ValidAt,
			Status:         item.Status,
			Category:       item.Category,
			Content:        item.Content,
			NeedUsers:      item.NeedUsers,
			UseTimes:       item.UseTimes,
			ExistingUsers:  item.ExistingUsers,
			AvailableTimes: item.AvailableTimes,
			CreatedAt:      item.CreatedAt,
		})
	}
	return demandList
}

func APIListDemands() []schema.APIDemandListResp {
	demands := mysql.GetAllDemandsByStatus(models.DemandStatusPublished)
	var demandList []schema.APIDemandListResp
	for idx := range demands {
		item := demands[idx]
		demandList = append(demandList, schema.APIDemandListResp{
			DemandId:  item.DemandId,
			Name:      item.Name,
			Brief:     item.Brief,
			ValidAt:   item.ValidAt,
			Category:  item.Category,
			CreatedAt: item.CreatedAt,
		})
	}
	return demandList
}

func APIListDemandContracts(category string) []schema.APIDemandContractListResp {
	demands := mysql.GetAllDemandsByStatusAndCategory(models.DemandStatusPublished, category)
	var demandList []schema.APIDemandContractListResp
	for idx := range demands {
		item := demands[idx]
		demandList = append(demandList, schema.APIDemandContractListResp{
			DemandId: item.DemandId,
			Name:     item.Name,
			Brief:    item.Brief,
			DemandContract: &schema.DemandContract{
				Address:   item.ContractAddr,
				ABI:       contracts.ApiMetaData.ABI,
				TokenName: item.ContractToken,
				Symbol:    item.ContractSymbol,
			},
		})
	}
	return demandList
}

func APIGetDemand(demandId int64) (*schema.APIDemandDetailResp, error) {
	demand, err := mysql.GetDemandDetail(demandId)
	if err != nil {
		return nil, err
	}
	return &schema.APIDemandDetailResp{
		DemandId: demand.DemandId,
		Name:     demand.Name,
		Brief:    demand.Brief,
		ValidAt:  demand.ValidAt,
		Category: demand.Category,
		Content:  demand.Content,
		DemandContract: &schema.DemandContract{
			Address:   demand.ContractAddr,
			ABI:       contracts.ApiMetaData.ABI,
			TokenName: demand.ContractToken,
			Symbol:    demand.ContractSymbol,
		},
		Purpose:   demand.Purpose,
		Agreement: demand.Agreement,
		CreatedAt: demand.CreatedAt,
	}, nil
}

func TotalDemands() int64 {
	return mysql.GetAllDemandsCount()
}

func TotalPublishedDemands() int64 {
	return mysql.GetAllDemandsByStatusCount(models.DemandStatusPublished)
}

func GetDemandDetail(demandId int64) (*schema.DemandDetailResp, error) {
	demand, err := mysql.GetDemandDetail(demandId)
	if err != nil {
		return nil, err
	}
	return &schema.DemandDetailResp{
		DemandId:        demand.DemandId,
		Name:            demand.Name,
		Brief:           demand.Brief,
		ValidAt:         demand.ValidAt,
		CreatedAt:       demand.CreatedAt,
		Category:        demand.Category,
		Content:         demand.Content,
		NeedUsers:       demand.NeedUsers,
		UseTimes:        demand.UseTimes,
		ExistingUsers:   demand.ExistingUsers,
		AvailableTimes:  demand.AvailableTimes,
		Purpose:         demand.Purpose,
		Algorithm:       demand.Algorithm,
		Agreement:       demand.Agreement,
		ContractAddress: demand.ContractAddr,
		ContractSymbol:  demand.ContractSymbol,
	}, nil
}

func GetDemandInfo(demandId int64) (*schema.DemandInfoResp, error) {
	demand, err := mysql.GetDemandDetail(demandId)
	if err != nil {
		return nil, err
	}
	return &schema.DemandInfoResp{
		DemandId:  demand.DemandId,
		Name:      demand.Name,
		Brief:     demand.Brief,
		ValidAt:   demand.ValidAt,
		Category:  demand.Category,
		Content:   demand.Content,
		NeedUsers: demand.NeedUsers,
		UseTimes:  demand.UseTimes,
		Purpose:   demand.Purpose,
		Algorithm: demand.Algorithm,
		Agreement: demand.Agreement,
	}, nil
}

func UpdateDemand(dur *schema.DemandUpdateReq) error {
	demand := models.Demand{
		DemandId:  dur.DemandId,
		Name:      dur.Name,
		Brief:     dur.Brief,
		ValidAt:   dur.ValidAt,
		Category:  dur.Category,
		Content:   dur.Content,
		NeedUsers: dur.NeedUsers,
		UseTimes:  dur.UseTimes,
		Purpose:   dur.Purpose,
		Algorithm: dur.Algorithm,
		Agreement: dur.Agreement,
	}
	return mysql.UpdateDemand(&demand)
}

func PublishDemand(demandId int64, client *ethclient.Client) error {
	// 1. 检查当前状态是否为草稿状态，草稿状态才可以发布
	if err := mysql.CheckDemandInitStatus(demandId); err != nil {
		return err
	}

	// 2. 部署合约
	privateKeyStr := viper.GetString("nft.goerli_private_key")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	result := make(chan error)
	go func() {
		// todo：read tokenName、symbol from db
		deployedContractAddress, deployTx, err := contracts.Deploy(client, privateKey, "PeopleDataBank", "PDB")
		if err != nil {
			result <- err
			return
		}
		zap.L().Info("contract deploy successfully", zap.String("address", deployedContractAddress), zap.String("transaction", deployTx))
		// 2.3 将合约地址更新到数据库，并设置状态为发布中...
		if err := mysql.UpdateDemandContract(demandId, deployedContractAddress, deployTx); err != nil {
			result <- err
			return
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case err := <-result:
		zap.L().Error("deploy contract error", zap.Error(err))
		return err
	}
}

func GetAllPublishingDemands() []models.Demand {
	return mysql.GetAllDemandsByStatus(models.DemandStatusPublishing)
}

// GetAllPublishedDemands 获取所有已经发布过的需求（已经有合约地址了）
func GetAllPublishedDemands() []models.Demand {
	return mysql.GetAllDemandsWithContracts(models.DemandStatusPublished)
}

// DemandStatusCronWorker 需求状态任务定时器
func DemandStatusCronWorker(client *ethclient.Client, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(demand.ContractTx))
	demandIdstr := strconv.FormatInt(demand.DemandId, 10)
	if err != nil {
		zap.L().Error("Demand contract status check cron error", zap.String("reason", "eth transaction receipt error"),
			zap.String("demand_id", demandIdstr),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("contract_tx", demand.ContractTx), zap.Error(err))
		return
	}
	// 执行成功，更新合约状态
	if receipt != nil {
		var err error
		if receipt.Status == types.ReceiptStatusSuccessful {
			err = mysql.UpdateDemandStatus(demand.DemandId, models.DemandStatusPublished)
		} else if receipt.Status == types.ReceiptStatusFailed {
			err = mysql.UpdateDemandStatus(demand.DemandId, models.DemandStatusPublishFailed)
		}
		if err != nil {
			zap.L().Error("Demand contract status check cron error", zap.String("reason", "update demand status error"),
				zap.String("demand_id", demandIdstr),
				zap.String("contract_address", demand.ContractAddr),
				zap.String("contract_tx", demand.ContractTx), zap.Error(err))
			return
		}
	}
}

// DemandContractRecordsCronWorker 获取签约数据记录
func DemandContractRecordsCronWorker(etherscanApiKey string, client *ethclient.Client, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	zap.L().Debug("Demand contract record starting", zap.String("ContractAddress", demand.ContractAddr))

	// 获取合约实例
	instanceApi, err := contracts.NewApi(common.HexToAddress(demand.ContractAddr), client)
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "get eth contract instance error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("contract_tx", demand.ContractTx), zap.Error(err))
		return
	}

	// 连接IPFS API的客户端
	ipfsClient := ipfsapi.NewShell("localhost:5001")
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "connect ipfs api error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr), zap.Error(err))
		return
	}

	// 查询该NFT的交易历史记录
	// todo：超过1000个分页
	url := fmt.Sprintf("https://api-goerli.etherscan.io/api?module=account&action=tokennfttx&contractaddress=%s&page=1&offset=1000&sort=desc&apikey=%s",
		demand.ContractAddr, etherscanApiKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "request tokennfttx error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("etherscan_url", url), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "decode request tokennfttx error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr), zap.Error(err))
		return
	}

	if result["status"].(string) == "1" {
		transactions := result["result"].([]interface{})
		for _, tx := range transactions {
			txHash := tx.(map[string]interface{})["hash"].(string)
			blockNumber := tx.(map[string]interface{})["blockNumber"].(string)
			timeStamp := tx.(map[string]interface{})["timeStamp"].(string)
			from := tx.(map[string]interface{})["from"].(string)
			to := tx.(map[string]interface{})["to"].(string)
			tokenID := tx.(map[string]interface{})["tokenID"].(string)

			fmt.Printf("TokenID: %s\n交易哈希：%s\n区块高度：%s\n交易时间戳：%s\n发送地址：%s\n接收地址：%s\n\n",
				tokenID, txHash, blockNumber, timeStamp, from, to)

			// 如果是发给合约地址的，则表示授权给合约地址的，可以直接处理该 tokenID
			if from != "0x0000000000000000000000000000000000000000" && to == strings.ToLower(demand.ContractAddr) {
				// 如果当前 token id 已经处理过存入数据库了，则忽略
				tokenIDInt, _ := strconv.ParseInt(tokenID, 10, 64)
				_, err := mysql.GetContractRecordByTokenId(demand.DemandId, tokenIDInt)
				if err != nil {
					// 如果不是不存在的错误，则记录日志后退出
					if !errors.Is(err, mysql.ErrContractRecordNotExist) {
						zap.L().Error("Demand contract records get cron error", zap.String("reason", "mysql.GetContractRecordByTokenId error"),
							zap.Int64("demand_id", demand.DemandId),
							zap.String("contract_address", demand.ContractAddr),
							zap.String("tokenId", tokenID), zap.Error(err))
						continue
					}
				} else {
					// 找到了数据，则不处理了
					continue
				}

				// 获取 tokenID 的 URI
				assetData, err := instanceApi.TokenURI(&bind.CallOpts{Context: ctx}, big.NewInt(tokenIDInt))
				if err != nil {
					zap.L().Error("Demand contract records get cron error", zap.String("reason", "get TokenURI error"),
						zap.Int64("demand_id", demand.DemandId),
						zap.String("contract_address", demand.ContractAddr),
						zap.String("tokenId", tokenID), zap.Error(err))
					continue
				}
				zap.L().Debug("get asset data by tokenId",
					zap.Int64("demand_id", demand.DemandId),
					zap.String("tokenId", tokenID),
					zap.String("tokenURI", assetData))

				// 解析 nft 中存放在 ipfs 上数据
				ipfsPath := path.New(assetData[7:])
				data, err := ipfs.Read(ipfsClient, ipfsPath.String())
				if err != nil {
					zap.L().Error("Demand contract records get cron error", zap.String("reason", "ipfs client read error"),
						zap.Int64("demand_id", demand.DemandId),
						zap.String("contract_address", demand.ContractAddr),
						zap.String("tokenId", tokenID),
						zap.String("tokenURI", assetData),
						zap.String("ipfsPath", ipfsPath.String()), zap.Error(err))
					continue
				}

				zap.L().Debug("get nft data successfully",
					zap.Int64("demand_id", demand.DemandId),
					zap.String("tokenId", tokenID),
					zap.String("tokenURI", assetData),
					zap.String("ipfsData", string(data)))

				var nftData schema.NftData
				if err := json.Unmarshal(data, &nftData); err != nil {
					zap.L().Error("Demand contract records get cron error", zap.String("reason", "json unmarshal ipfs data error"),
						zap.Int64("demand_id", demand.DemandId),
						zap.String("contract_address", demand.ContractAddr),
						zap.String("tokenId", tokenID),
						zap.String("tokenURI", assetData),
						zap.String("ipfsPath", ipfsPath.String()), zap.String("ipfsData", string(data)), zap.Error(err))
					continue
				}

				// 解析did doc数据
				didDocPath := path.New(nftData.DidDoc[7:])
				// read ipfs data
				didDocData, err := ipfs.Read(ipfsClient, didDocPath.String())
				if err != nil {
					zap.L().Error("Demand contract records get cron error", zap.String("reason", "ipfs read did doc error"),
						zap.Int64("demand_id", demand.DemandId),
						zap.String("contract_address", demand.ContractAddr),
						zap.String("tokenId", tokenID),
						zap.String("tokenURI", assetData),
						zap.String("didDocPath", didDocPath.String()), zap.Error(err))
					continue
				}

				zap.L().Debug("get did doc data successfully",
					zap.Int64("demand_id", demand.DemandId),
					zap.String("tokenId", tokenID),
					zap.String("tokenURI", assetData),
					zap.String("ipfsData", string(data)), zap.String("didDocPath", didDocPath.String()))

				r := regexp.MustCompile(`"did":"(.*?)"`)
				match := r.FindStringSubmatch(string(didDocData))
				if len(match) > 1 {
					seconds, _ := strconv.ParseInt(timeStamp, 10, 64)
					record := models.ContractRecord{
						DemandId: demand.DemandId,
						TokenId:  tokenIDInt,
						TokenURI: assetData,
						DidURI:   nftData.DidDoc,
						Did:      match[1],
						SignTime: time.Unix(seconds, 0),
					}
					err := mysql.InsertContractRecord(&record)
					if err != nil {
						zap.L().Error("Demand contract records get cron error", zap.String("reason", "insert contract record error"),
							zap.Int64("demand_id", demand.DemandId),
							zap.String("contract_address", demand.ContractAddr),
							zap.String("tokenId", tokenID),
							zap.String("tokenURI", assetData),
							zap.String("did", match[1]),
							zap.String("didDocPath", didDocPath.String()), zap.Error(err))
						continue
					}
					zap.L().Debug("insert contract record data successfully",
						zap.Int64("demand_id", demand.DemandId),
						zap.String("tokenId", tokenID),
						zap.String("tokenURI", assetData),
						zap.String("ipfsData", string(data)),
						zap.String("did", match[1]),
						zap.String("didDocPath", didDocPath.String()))
				}
			}
		}
	} else {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "request tokennfttx failed"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("message", result["message"].(string)))
	}
}

// DemandContractRecordsCronWorker 获取签约数据记录
func DemandContractRecordsCron(client *ethclient.Client, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Get a contract instance
	contractAddr := common.HexToAddress(demand.ContractAddr) // 合约地址
	instanceApi, err := contracts.NewApi(contractAddr, client)
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "get eth contract instance error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr),
			zap.String("contract_tx", demand.ContractTx), zap.Error(err))
		return
	}

	// Get the total number of tokens
	totalTokens, err := instanceApi.TotalSupply(&bind.CallOpts{Context: ctx})
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "get TotalSupply error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr), zap.Error(err))
		return
	}

	zap.L().Debug("get token supply of contract", zap.String("contract", demand.ContractAddr), zap.String("totalTokens", totalTokens.String()))

	// Connect to the IPFS API
	ipfsClient := ipfsapi.NewShell("localhost:5001")
	if err != nil {
		zap.L().Error("Demand contract records get cron error", zap.String("reason", "connect ipfs api error"),
			zap.Int64("demand_id", demand.DemandId),
			zap.String("contract_address", demand.ContractAddr), zap.Error(err))
		return
	}

	// Loop over all token IDs and retrieve the owner of each token
	for i := big.NewInt(0); i.Cmp(totalTokens) < 0; {
		i = i.Add(i, big.NewInt(1))

		// 如果当前 token id 已经处理过存入数据库了，则忽略
		_, err := mysql.GetContractRecordByTokenId(demand.DemandId, i.Int64())
		if err != nil {
			// 如果不是不存在的错误，则记录日志后退出
			if !errors.Is(err, mysql.ErrContractRecordNotExist) {
				zap.L().Error("Demand contract records get cron error", zap.String("reason", "mysql.GetContractRecordByTokenId error"),
					zap.Int64("demand_id", demand.DemandId),
					zap.String("contract_address", demand.ContractAddr),
					zap.String("tokenId", i.String()), zap.Error(err))
				continue
			}
		} else {
			// 找到了数据，则不处理了
			continue
		}

		// 获取当前 token id 的 owner
		owner, err := instanceApi.OwnerOf(&bind.CallOpts{Context: ctx}, i)
		if err != nil {
			zap.L().Error("Demand contract records get cron error", zap.String("reason", "get OwnerOf error"),
				zap.Int64("demand_id", demand.DemandId),
				zap.String("contract_address", demand.ContractAddr),
				zap.String("tokenId", i.String()), zap.Error(err))
			continue
		}

		zap.L().Debug("get owned by tokenId", zap.Int64("demand_id", demand.DemandId),
			zap.String("tokenId", i.String()), zap.String("owner", owner.String()))

		// 如果owner是当前合约，则表面是用户将nft授权转移给了合约
		if owner.String() == demand.ContractAddr {
			assetData, err := instanceApi.TokenURI(&bind.CallOpts{Context: ctx}, i)
			if err != nil {
				zap.L().Error("Demand contract records get cron error", zap.String("reason", "get TokenURI error"),
					zap.Int64("demand_id", demand.DemandId),
					zap.String("contract_address", demand.ContractAddr),
					zap.String("tokenId", i.String()), zap.Error(err))
				continue
			}
			zap.L().Debug("get asset data by tokenId",
				zap.Int64("demand_id", demand.DemandId),
				zap.String("tokenId", i.String()),
				zap.String("tokenURI", assetData))

			// Parse the IPFS path
			ipfsPath := path.New(assetData[7:])
			// read ipfs data
			data, err := ipfs.Read(ipfsClient, ipfsPath.String())
			if err != nil {
				zap.L().Error("Demand contract records get cron error", zap.String("reason", "ipfs client read error"),
					zap.Int64("demand_id", demand.DemandId),
					zap.String("contract_address", demand.ContractAddr),
					zap.String("tokenId", i.String()),
					zap.String("tokenURI", assetData),
					zap.String("ipfsPath", ipfsPath.String()), zap.Error(err))
				continue
			}

			zap.L().Debug("get nft data successfully",
				zap.Int64("demand_id", demand.DemandId),
				zap.String("tokenId", i.String()),
				zap.String("tokenURI", assetData),
				zap.String("ipfsData", string(data)))

			var nftData schema.NftData
			if err := json.Unmarshal(data, &nftData); err != nil {
				zap.L().Error("Demand contract records get cron error", zap.String("reason", "json unmarshal ipfs data error"),
					zap.Int64("demand_id", demand.DemandId),
					zap.String("contract_address", demand.ContractAddr),
					zap.String("tokenId", i.String()),
					zap.String("tokenURI", assetData),
					zap.String("ipfsPath", ipfsPath.String()), zap.String("ipfsData", string(data)), zap.Error(err))
				continue
			}

			// Parse the IPFS path
			didDocPath := path.New(nftData.DidDoc[7:])
			// read ipfs data
			didDocData, err := ipfs.Read(ipfsClient, didDocPath.String())
			if err != nil {
				zap.L().Error("Demand contract records get cron error", zap.String("reason", "ipfs read did doc error"),
					zap.Int64("demand_id", demand.DemandId),
					zap.String("contract_address", demand.ContractAddr),
					zap.String("tokenId", i.String()),
					zap.String("tokenURI", assetData),
					zap.String("didDocPath", didDocPath.String()), zap.Error(err))
				continue
			}

			zap.L().Debug("get did doc data successfully",
				zap.Int64("demand_id", demand.DemandId),
				zap.String("tokenId", i.String()),
				zap.String("tokenURI", assetData),
				zap.String("ipfsData", string(data)), zap.String("didDocPath", didDocPath.String()))

			r := regexp.MustCompile(`"did":"(.*?)"`)
			match := r.FindStringSubmatch(string(didDocData))
			if len(match) > 1 {
				record := models.ContractRecord{
					DemandId: demand.DemandId,
					TokenId:  i.Int64(),
					TokenURI: assetData,
					DidURI:   nftData.DidDoc,
					Did:      match[1],
					//SignTime: time.Now(), // todo：查找正确的签约时间
				}
				err := mysql.InsertContractRecord(&record)
				if err != nil {
					zap.L().Error("Demand contract records get cron error", zap.String("reason", "insert contract record error"),
						zap.Int64("demand_id", demand.DemandId),
						zap.String("contract_address", demand.ContractAddr),
						zap.String("tokenId", i.String()), zap.String("tokenURI", assetData),
						zap.String("did", match[1]),
						zap.String("didDocPath", didDocPath.String()), zap.Error(err))
					continue
				}
				zap.L().Debug("insert contract record data successfully",
					zap.Int64("demand_id", demand.DemandId),
					zap.String("tokenId", i.String()),
					zap.String("tokenURI", assetData),
					zap.String("ipfsData", string(data)),
					zap.String("did", match[1]),
					zap.String("didDocPath", didDocPath.String()))
			} else {
				zap.L().Debug("get did doc data successfully but no match did data",
					zap.Int64("demand_id", demand.DemandId),
					zap.String("tokenId", i.String()),
					zap.String("tokenURI", assetData),
					zap.String("ipfsData", string(data)),
					zap.String("didDocPath", didDocPath.String()))
			}

		}

	}

}
