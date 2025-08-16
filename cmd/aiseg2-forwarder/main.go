package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/aiseg2"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/config"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/influx"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/pkg/logger"
)

func main() {
	// ロガー初期化
	logger.Setup()

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		slog.Error("設定の読み込みに失敗しました", "error", err)
		os.Exit(1)
	}

	// 設定を表示
	slog.Info("設定情報",
		"aiseg2_host", cfg.AiSEG2.Host,
		"aiseg2_user", cfg.AiSEG2.User,
		"aiseg2_use_https", cfg.AiSEG2.UseHTTPS,
		"influxdb_host", cfg.InfluxDB.Host,
		"influxdb_org", cfg.InfluxDB.Org,
		"influxdb_bucket", cfg.InfluxDB.Bucket,
		"influxdb_use_https", cfg.InfluxDB.UseHTTPS,
	)

	// AiSEG2クライアント初期化
	aiseg2Client, err := aiseg2.NewClient(cfg.AiSEG2)
	if err != nil {
		slog.Error("AiSEG2クライアントの初期化に失敗しました", "error", err)
		os.Exit(1)
	}

	// InfluxDBクライアント初期化
	influxClient, err := influx.NewClient(cfg.InfluxDB)
	if err != nil {
		slog.Error("InfluxDBクライアントの初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	defer influxClient.Close()

	// コンテキストとシグナルハンドリング
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// メインループをgoroutineで実行
	go runMainLoop(ctx, aiseg2Client, influxClient)

	// シグナル待機
	<-sigChan
	slog.Info("シャットダウンシグナルを受信しました")
	cancel()

	// 少し待ってからプロセス終了
	time.Sleep(1 * time.Second)
	slog.Info("アプリケーションを終了します")
}

func runMainLoop(ctx context.Context, aiseg2Client *aiseg2.Client, influxClient *influx.Client) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("メインループを終了します")
			return
		case <-ticker.C:
			if err := executeDataCollection(aiseg2Client, influxClient); err != nil {
				slog.Error("データ収集でエラーが発生しました", "error", err)
			}
		}
	}
}

func executeDataCollection(aiseg2Client *aiseg2.Client, influxClient *influx.Client) error {
	now := time.Now()

	// AiSEG2からデータを取得
	powerSummary, err := aiseg2Client.GetPowerSummary()
	if err != nil {
		return err
	}
	slog.Info("電力サマリーを取得しました",
		"timestamp", now.Format("2006-01-02 15:04:05"),
		"data", powerSummary,
	)

	detailsUsagePower, err := aiseg2Client.GetDetailsUsagePower()
	if err != nil {
		return err
	}
	slog.Info("詳細消費電力を取得しました",
		"timestamp", now.Format("2006-01-02 15:04:05"),
		"count", len(detailsUsagePower),
	)

	// InfluxDBへデータを送信
	if err := influxClient.WritePower(powerSummary, detailsUsagePower); err != nil {
		return err
	}

	slog.Info("データ書き込み完了")
	return nil
}
