package logic

import (
	"context"
	"patronus/internal/dao/mysql"
	"patronus/internal/models"
	"patronus/internal/schema"
	"patronus/pkg/contracts"
	"patronus/pkg/snowflake"
	"strconv"
	"sync"
	"time"

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

func DemandCronWorker(client *ethclient.Client, demand *models.Demand, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(demand.ContractTx))
	demandIdstr := strconv.FormatInt(demand.DemandId, 10)
	if err != nil {
		zap.L().Error("Demand contract status check cron error", zap.String("reason", "etch transaction receipt error"),
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
