package api

import (
	"errors"
	"patronus/internal/controller"
	"patronus/internal/dao/mysql"
	"patronus/internal/logic"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func DemandListHandler(c *gin.Context) {
	// todo：分页
	controller.ResponseOk(c, gin.H{
		"demands": logic.APIListDemands(),
		"total":   logic.TotalPublishedDemands(),
	})

}

func DemandContractListHandler(c *gin.Context) {
	// 1. 获取demand id
	categoryStr := c.Param("category")
	controller.ResponseOk(c, gin.H{
		"demands": logic.APIListDemandContracts(categoryStr),
	})
}

func DemandDetailHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("DemandDetail API Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	data, err := logic.APIGetDemand(demandId)
	if err != nil {
		zap.L().Error("DemandDetail API Handler get detail error", zap.String("id", demandIdStr), zap.Error(err))
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
