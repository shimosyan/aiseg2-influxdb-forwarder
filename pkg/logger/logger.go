package logger

import (
	"log/slog"
	"os"
)

// Setup はログ設定を初期化します
func Setup() {
	// JSON形式でのログ出力（コンテナ環境に適している）
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
