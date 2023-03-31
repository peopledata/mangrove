package admin

import (
	controller "mangrove/internal/controller"
	"mangrove/internal/logic"
	"mangrove/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TaskCreateHandler(c *gin.Context) {
	// 1. 获取到当前操作用户
	val, ok := c.Get(controller.CtxUserKey)
	if !ok {
		controller.ResponseErr(c, controller.CodeUnknown)
		return
	}
	user, _ := val.(*models.User)

	//	2.参数校验
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("TaskCreate Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 3. 业务处理
	if err := logic.CreateTask(demandId, user); err != nil {
		zap.L().Error("TaskCreate Handler logic handle error", zap.String("demand_id",
			strconv.FormatInt(demandId, 10)), zap.Error(err))
		controller.ResponseErr(c, controller.CodeTaskCreateErr)
		return
	}

	// 4. 返回Response
	controller.ResponseOk(c, nil)
}

func TaskListHandler(c *gin.Context) {
	// 1. 获取demand id
	demandIdStr := c.Param("id")
	demandId, err := strconv.ParseInt(demandIdStr, 10, 64)
	if err != nil {
		zap.L().Error("TaskList Handler with invalid param", zap.String("id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	tasks, err := logic.ListTasks(demandId)
	if err != nil {
		zap.L().Error("TaskList Handler logic handle error", zap.String("demand_id", demandIdStr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeTaskListErr)
		return
	}

	controller.ResponseOk(c, tasks)
}

func TaskAlgoHandler(c *gin.Context) {
	// 1. 获取demand id
	taskIdstr := c.Param("id")
	taskId, err := strconv.ParseInt(taskIdstr, 10, 64)
	if err != nil {
		zap.L().Error("Task Algo API Handler with invalid param", zap.String("task_id", taskIdstr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeInvalidParam)
		return
	}

	// 2. 业务处理
	records, err := logic.TaskAlgoRecords(taskId)
	if err != nil {
		zap.L().Error("Task Algo API Handler logic handle error", zap.String("task_id", taskIdstr), zap.Error(err))
		controller.ResponseErr(c, controller.CodeAlgoRecordsErr)
		return
	}

	controller.ResponseOk(c, gin.H{
		"records": records,
	})
}
