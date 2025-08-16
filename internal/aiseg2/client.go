package aiseg2

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shimosyan/aiseg2-influxdb-forwarder/internal/config"
	dac "github.com/xinsnake/go-http-digest-auth-client"
)

// Client はAiSEG2 WebUIとの通信を管理します
type Client struct {
	host     string
	useHTTPS bool
	client   *http.Client
}

// NewClient は新しいAiSEG2クライアントを作成します
func NewClient(cfg config.AiSEG2Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("AiSEG2のホストが指定されていません")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("AiSEG2のログインユーザー名が指定されていません")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("AiSEG2のログインパスワードが指定されていません")
	}

	// Digest認証クライアントの作成
	transport := dac.NewTransport(cfg.User, cfg.Password)
	client := &http.Client{Transport: &transport}

	return &Client{
		host:     cfg.Host,
		useHTTPS: cfg.UseHTTPS,
		client:   client,
	}, nil
}

// getNumericValue は文字列から数値を抽出します
func (c *Client) getNumericValue(input string) float64 {
	if input == "" {
		return 0
	}

	// 数字とピリオドのみを抽出
	re := regexp.MustCompile(`[0-9.]+`)
	matches := re.FindAllString(input, -1)
	if len(matches) == 0 {
		return 0
	}

	// 最初にマッチした数値を返す
	value, err := strconv.ParseFloat(matches[0], 64)
	if err != nil {
		return 0
	}
	return value
}

// getURL はスキームとホストを組み合わせてURLを生成します
func (c *Client) getURL(path string) string {
	scheme := "http"
	if c.useHTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, c.host, path)
}

// GetPowerSummary は電力サマリー情報を取得します
func (c *Client) GetPowerSummary() (*PowerSummary, error) {
	url := c.getURL("/page/electricflow/111")

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("AiSEG2への接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AiSEG2からの応答が異常です: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	// 総発電電力
	totalGenerationPowerKW := MetricsElement{
		Name:  "総発電電力(kW)",
		Value: c.getNumericValue(doc.Find("#g_capacity").Text()),
	}

	// 総消費電力
	totalUsagePowerKW := MetricsElement{
		Name:  "総消費電力(kW)",
		Value: c.getNumericValue(doc.Find("#u_capacity").Text()),
	}

	// 売買電力（計算値）
	totalBalancePowerKW := MetricsElement{
		Name:  "売買電力(kW)",
		Value: totalGenerationPowerKW.Value - totalUsagePowerKW.Value,
	}

	// 詳細発電電力
	var detailsGenerationPower []MetricsElement
	for i := 1; i <= 3; i++ {
		titleSelector := fmt.Sprintf("#g_d_%d_title", i)
		capacitySelector := fmt.Sprintf("#g_d_%d_capacity", i)

		title := strings.TrimSpace(doc.Find(titleSelector).Text())
		if title != "" {
			detailsGenerationPower = append(detailsGenerationPower, MetricsElement{
				Name:  fmt.Sprintf("%s(W)", title),
				Value: c.getNumericValue(doc.Find(capacitySelector).Text()),
			})
		}
	}

	return &PowerSummary{
		TotalGenerationPowerKW:  totalGenerationPowerKW,
		TotalUsagePowerKW:       totalUsagePowerKW,
		TotalBalancePowerKW:     totalBalancePowerKW,
		DetailsGenerationPower:  detailsGenerationPower,
	}, nil
}

// GetDetailsUsagePower は詳細消費電力情報を取得します
func (c *Client) GetDetailsUsagePower() (DetailUsagePower, error) {
	var usagePowerItems []MetricsElement
	pageEndCheck := ""
	maxCount := 20

	for pageCount := 1; pageCount <= maxCount; pageCount++ {
		url := c.getURL(fmt.Sprintf("/page/electricflow/1113?id=%d", pageCount))
		
		resp, err := c.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("AiSEG2への接続に失敗しました: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("AiSEG2からの応答が異常です: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("HTMLの解析に失敗しました: %w", err)
		}

		// 重複ページチェック用
		var checkDuplicate []string
		for i := 1; i <= 10; i++ {
			deviceSelector := fmt.Sprintf("#stage_%d > div.c_device", i)
			name := strings.TrimSpace(doc.Find(deviceSelector).Text())
			checkDuplicate = append(checkDuplicate, name)
		}

		currentPageCheck := strings.Join(checkDuplicate, ",")
		if pageEndCheck == currentPageCheck {
			break // 重複ページに到達
		}

		// データ抽出
		for i := 1; i <= 10; i++ {
			deviceSelector := fmt.Sprintf("#stage_%d > div.c_device", i)
			valueSelector := fmt.Sprintf("#stage_%d > div.c_value", i)
			
			name := strings.TrimSpace(doc.Find(deviceSelector).Text())
			if name == "" {
				continue
			}

			usagePowerItems = append(usagePowerItems, MetricsElement{
				Name:  fmt.Sprintf("%s(W)", name),
				Value: c.getNumericValue(doc.Find(valueSelector).Text()),
			})
		}

		pageEndCheck = currentPageCheck
	}

	return usagePowerItems, nil
}
