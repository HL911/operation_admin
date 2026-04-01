package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"operation_admin/backend/internal/config"
	"operation_admin/backend/internal/database"
	"operation_admin/backend/internal/http/handler"
	"operation_admin/backend/internal/http/router"
	"operation_admin/backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Application 管理应用初始化结果与整个服务生命周期。
type Application struct {
	// Config 保存当前进程生效的完整运行配置，供启动日志和下游组件复用。
	Config *config.Config
	// Logger 提供应用级结构化日志能力，覆盖启动、关闭和运行期关键事件。
	Logger *zap.Logger
	// DB 持有 PostgreSQL 对应的 Gorm 连接实例，供业务层后续复用。
	DB *gorm.DB
	// HTTPServer 持有标准库 HTTP Server 实例，负责真正的端口监听与请求处理。
	HTTPServer *http.Server
	// ShutdownTimeout 指定优雅停机允许的最长等待时间，单位为 Go duration。
	ShutdownTimeout time.Duration
}

// New 负责装配配置、日志、数据库和 HTTP 服务。
func New(configPath string) (*Application, error) {
	// loadedConfig 保存从文件和环境变量合并后的运行配置。
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// zapLogger 提供全局结构化日志能力。
	zapLogger, err := logger.New(loadedConfig.Log)
	if err != nil {
		return nil, fmt.Errorf("初始化日志失败: %w", err)
	}

	// gormDB 维护 PostgreSQL 连接池与 ORM 能力。
	gormDB, err := database.New(loadedConfig.Database)
	if err != nil {
		_ = logger.Sync(zapLogger)
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	// healthHandler 提供服务存活检查接口。
	healthHandler := handler.NewHealthHandler(loadedConfig.App.Name)

	// engine 负责承载 Gin 路由与中间件链。
	engine := router.New(loadedConfig.Server, zapLogger, healthHandler)

	// httpServer 是实际对外监听端口的标准库服务对象。
	httpServer := &http.Server{
		Addr:              loadedConfig.Server.Address(),
		Handler:           engine,
		ReadHeaderTimeout: loadedConfig.Server.ReadTimeout,
		ReadTimeout:       loadedConfig.Server.ReadTimeout,
		WriteTimeout:      loadedConfig.Server.WriteTimeout,
	}

	return &Application{
		Config:          loadedConfig,
		Logger:          zapLogger,
		DB:              gormDB,
		HTTPServer:      httpServer,
		ShutdownTimeout: loadedConfig.Server.ShutdownTimeout,
	}, nil
}

// Run 负责启动服务并在收到系统信号后执行优雅关闭。
func (a *Application) Run() error {
	// signalContext 用于监听进程级退出信号。
	signalContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// serverErrorCh 用于接收 HTTP 服务启动或运行期间返回的错误。
	serverErrorCh := make(chan error, 1)

	a.Logger.Info(
		"应用启动完成，准备监听请求",
		zap.String("app_name", a.Config.App.Name),
		zap.String("address", a.HTTPServer.Addr),
		zap.String("mode", a.Config.Server.Mode),
	)

	go func() {
		serverErrorCh <- a.startHTTPServer()
	}()

	select {
	case err := <-serverErrorCh:
		return err
	case <-signalContext.Done():
		a.Logger.Info("收到停止信号，开始优雅关闭")
		return a.shutdownHTTPServer()
	}
}

// Close 负责释放数据库连接与日志资源。
func (a *Application) Close() error {
	// closeErrors 汇总资源释放阶段出现的错误。
	var closeErrors []error

	if a.DB != nil {
		if err := database.Close(a.DB); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("关闭数据库失败: %w", err))
		}
	}

	if a.Logger != nil {
		if err := logger.Sync(a.Logger); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("同步日志失败: %w", err))
		}
	}

	if len(closeErrors) == 0 {
		return nil
	}

	return errors.Join(closeErrors...)
}

// startHTTPServer 负责启动 HTTP 服务并记录监听信息。
func (a *Application) startHTTPServer() error {
	a.Logger.Info("HTTP 服务开始监听", zap.String("address", a.HTTPServer.Addr))

	if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// shutdownHTTPServer 负责基于默认超时时间执行优雅关闭。
func (a *Application) shutdownHTTPServer() error {
	// shutdownContext 用于限制服务优雅关闭的最长等待时间。
	shutdownContext, cancel := context.WithTimeout(context.Background(), a.ShutdownTimeout)
	defer cancel()

	a.Logger.Info("HTTP 服务开始关闭")
	return a.HTTPServer.Shutdown(shutdownContext)
}
