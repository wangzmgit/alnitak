package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"interastral-peace.com/alnitak/internal/config"
	"interastral-peace.com/alnitak/internal/worker"
	"interastral-peace.com/alnitak/pkg/logger"
)

var (
	env         = flag.String("env", "prod", "运行环境 dev/prod")
	concurrency = flag.Int("concurrency", 0, "最大并发转码数 (0=使用配置文件的 worker_concurrency)")
	workerID    = flag.String("id", "", "Worker 唯一标识 (默认取 hostname)")
	healthPort  = flag.Int("health-port", 9100, "健康检查 HTTP 端口")
)

func main() {
	flag.Parse()

	// 加载配置
	viper.SetConfigName(getConfigName(*env))
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: 配置文件读取失败: %v\n", err)
		panic("配置文件读取失败: " + err.Error())
	}

	cfg := &config.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: 配置解析失败: %v\n", err)
		panic("配置解析失败: " + err.Error())
	}

	// 解析并发数：优先 CLI 传入的 --concurrency，否则从配置文件读取
	if *concurrency <= 0 {
		if cfg.Transcoding.WorkerConcurrency > 0 {
			*concurrency = cfg.Transcoding.WorkerConcurrency
		} else {
			*concurrency = 2 // 兜底默认值
		}
	}

	// 初始化日志
	logger.InitLogger(cfg)
	zap.L().Info("transcoder-worker starting",
		zap.String("env", *env),
		zap.Int("concurrency", *concurrency),
	)

	if *workerID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			*workerID = fmt.Sprintf("unknown-%d", os.Getpid())
		} else {
			*workerID = hostname
		}
	}

	// 创建 Worker
	w, err := worker.NewWorker(cfg, *workerID, *concurrency)
	if err != nil {
		zap.L().Fatal("create worker failed", zap.Error(err))
	}

	// 信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		zap.L().Info("shutting down...")
		cancel()
	}()

	// 健康检查 HTTP server
	healthSrv := startHealthServer(w, *healthPort)
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		healthSrv.Shutdown(shutdownCtx)
	}()

	// 启动消费循环
	if err := w.Run(ctx); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "FATAL: worker run failed: %v\n", err)
		zap.L().Fatal("worker run failed", zap.Error(err))
	}

	zap.L().Info("transcoder-worker stopped")
}

// startHealthServer 启动健康检查 HTTP 服务。
func startHealthServer(w *worker.Worker, port int) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/ready", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		if w.Ready() {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"status": "ready",
				"health": w.Health(),
			})
		} else {
			rw.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"status": "not ready",
				"health": w.Health(),
			})
		}
	})

	mux.HandleFunc("/stats", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(w.Health())
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		zap.L().Info("Health server listening", zap.Int("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "FATAL: health server failed: %v\n", err)
				zap.L().Fatal("health server failed", zap.Error(err))
			}
	}()

	return srv
}

func getConfigName(env string) string {
	switch env {
	case "dev":
		return "application.dev"
	default:
		return "application.prod"
	}
}
