package influx

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/aiseg2"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/config"
)

// Client はInfluxDBとの通信を管理します
type Client struct {
	client influxdb2.Client
	bucket string
	org    string
}

// NewClient は新しいInfluxDBクライアントを作成します
func NewClient(cfg config.InfluxDBConfig) (*Client, error) {
	scheme := "http"
	if cfg.UseHTTPS {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s", scheme, cfg.Host)

	client := influxdb2.NewClient(url, cfg.Token)

	return &Client{
		client: client,
		bucket: cfg.Bucket,
		org:    cfg.Org,
	}, nil
}

// Close はInfluxDBクライアントを終了します
func (c *Client) Close() {
	c.client.Close()
}

// WritePower は電力データをInfluxDBに書き込みます
func (c *Client) WritePower(powerSummary *aiseg2.PowerSummary, detailsUsagePower aiseg2.DetailUsagePower) error {
	writeAPI := c.client.WriteAPIBlocking(c.org, c.bucket)

	var points []*write.Point
	now := time.Now()

	// 総発電電力
	points = append(points, influxdb2.NewPoint("power",
		map[string]string{"summary": powerSummary.TotalGenerationPowerKW.Name},
		map[string]interface{}{"value": powerSummary.TotalGenerationPowerKW.Value},
		now))

	// 総消費電力
	points = append(points, influxdb2.NewPoint("power",
		map[string]string{"summary": powerSummary.TotalUsagePowerKW.Name},
		map[string]interface{}{"value": powerSummary.TotalUsagePowerKW.Value},
		now))

	// 売買電力
	points = append(points, influxdb2.NewPoint("power",
		map[string]string{"summary": powerSummary.TotalBalancePowerKW.Name},
		map[string]interface{}{"value": powerSummary.TotalBalancePowerKW.Value},
		now))

	// 詳細発電電力
	for _, item := range powerSummary.DetailsGenerationPower {
		points = append(points, influxdb2.NewPoint("power",
			map[string]string{
				"detail-type":    "generation",
				"detail-section": item.Name,
			},
			map[string]interface{}{"value": item.Value},
			now))
	}

	// 詳細消費電力
	for _, item := range detailsUsagePower {
		points = append(points, influxdb2.NewPoint("power",
			map[string]string{
				"detail-type":    "usage",
				"detail-section": item.Name,
			},
			map[string]interface{}{"value": item.Value},
			now))
	}

	// バッチ書き込み実行
	err := writeAPI.WritePoint(context.Background(), points...)
	if err != nil {
		return fmt.Errorf("InfluxDBへの書き込みに失敗しました: %w", err)
	}

	return nil
}
