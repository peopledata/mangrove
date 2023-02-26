package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"patronus/internal/controller/admin"
	mysql2 "patronus/internal/dao/mysql"
	"patronus/internal/logger"
	"patronus/internal/routes"
	"patronus/pkg/converter"
	"patronus/pkg/snowflake"
	"patronus/settings"
	"syscall"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	//	1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}

	//	2. 初始化日志
	if err := logger.Init(viper.GetString("app.mode")); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	//	3. 初始化MySQL
	if err := mysql2.Init(); err != nil {
		fmt.Printf("init mysql failed, err: %v\n", err)
		return
	}
	defer mysql2.Close()

	mysql2.AutoMigrate()

	//	4.初始化Redis
	//if err := redis.Init(); err != nil {
	//	fmt.Printf("init redis failed, err: %v\n", err)
	//	return
	//}
	//defer redis.Close()

	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}

	// 定时器任务
	c := cron.New()
	err := c.AddFunc("*/5 * * * *", admin.DemandContractStatusCron)
	if err != nil {
		fmt.Printf("add cron task failed, err: %v\n", err)
		return
	}
	c.Start()

	// 初始化gin框架内置的validator使用的翻译器
	if err := converter.InitTrans("zh"); err != nil {
		fmt.Printf("init gin validator translate, err: %v\n", err)
		return
	}

	//	5. 注册路由
	r := routes.Setup(viper.GetString("app.mode"))

	//	6.启动服务（优雅退出）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
