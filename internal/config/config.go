package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	AiSEG2   AiSEG2Config
	InfluxDB InfluxDBConfig
}

// AiSEG2Config はAiSEG2関連の設定を保持します
type AiSEG2Config struct {
	Host     string
	User     string
	Password string
	UseHTTPS bool
}

// InfluxDBConfig はInfluxDB関連の設定を保持します
type InfluxDBConfig struct {
	Host     string
	Token    string
	Org      string
	Bucket   string
	UseHTTPS bool
}

// Load は環境変数から設定を読み込みます
func Load() (*Config, error) {
	// .envファイルを読み込み（ファイルが存在しない場合はエラーを無視）
	if err := godotenv.Load(); err != nil {
		// .envファイルが存在しない場合は警告を出力するが、エラーにはしない
		fmt.Printf("警告: .envファイルの読み込みに失敗しました（環境変数から読み込みます）: %v\n", err)
	}

	cfg := &Config{
		AiSEG2: AiSEG2Config{
			Host:     getEnvRequired("AISEG2_HOST"),
			User:     getEnvRequired("AISEG2_USER"),
			Password: getEnvRequired("AISEG2_PASSWORD"),
			UseHTTPS: getEnvBool("AISEG2_USE_HTTPS"),
		},
		InfluxDB: InfluxDBConfig{
			Host:     getEnvRequired("INFLUXDB_HOST"),
			Token:    getEnvRequired("INFLUXDB_TOKEN"),
			Org:      getEnvRequired("INFLUXDB_ORG"),
			Bucket:   getEnvRequired("INFLUXDB_BUCKET"),
			UseHTTPS: getEnvBool("INFLUXDB_USE_HTTPS"),
		},
	}

	return cfg, nil
}

// getEnvRequired は必須環境変数を取得します
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("環境変数 %s が設定されていません", key))
	}
	return value
}

// getEnvBool は環境変数をbool値として取得します
func getEnvBool(key string) bool {
	return os.Getenv(key) == "1"
}
